#include <stdlib.h>
#include <string.h>
#include <errno.h>
#include <limits.h>
#include <libappindicator/app-indicator.h>
#include "systray.h"

static AppIndicator *global_app_indicator;
static GtkWidget *global_tray_menu = NULL;
static GList *global_menu_items = NULL;
static char temp_file_name[PATH_MAX] = "";

typedef struct {
	GtkWidget *menu_item;
	int menu_id;
	long signalHandlerId;
} MenuItemNode;

typedef struct {
	int menu_id;
	int parent_menu_id;
	char* title;
	char* tooltip;
	short disabled;
	short checked;
	short isCheckable;
} MenuItemInfo;

void registerSystray(void) {
	gtk_init(0, NULL);
	global_app_indicator = app_indicator_new("systray", "", APP_INDICATOR_CATEGORY_APPLICATION_STATUS);
	app_indicator_set_status(global_app_indicator, APP_INDICATOR_STATUS_ACTIVE);
	global_tray_menu = gtk_menu_new();
	app_indicator_set_menu(global_app_indicator, GTK_MENU(global_tray_menu));
	systray_ready();
}

int nativeLoop(void) {
	gtk_main();
	systray_on_exit();
	return 0;
}

void _unlink_temp_file() {
	if (strlen(temp_file_name) != 0) {
		int ret = unlink(temp_file_name);
		if (ret == -1) {
			printf("failed to remove temp icon file %s: %s\n", temp_file_name, strerror(errno));
		}
		temp_file_name[0] = '\0';
	}
}

// runs in main thread, should always return FALSE to prevent gtk to execute it again
gboolean do_set_icon(gpointer data) {
	_unlink_temp_file();
	char *tmpdir = getenv("TMPDIR");
	if (NULL == tmpdir) {
		tmpdir = "/tmp";
	}
	strncpy(temp_file_name, tmpdir, PATH_MAX-1);
	strncat(temp_file_name, "/systray_XXXXXX", PATH_MAX-1);
	temp_file_name[PATH_MAX-1] = '\0';

	GBytes* bytes = (GBytes*)data;
	int fd = mkstemp(temp_file_name);
	if (fd == -1) {
		printf("failed to create temp icon file %s: %s\n", temp_file_name, strerror(errno));
		return FALSE;
	}
	gsize size = 0;
	gconstpointer icon_data = g_bytes_get_data(bytes, &size);
	ssize_t written = write(fd, icon_data, size);
	close(fd);
	if(written != size) {
		printf("failed to write temp icon file %s: %s\n", temp_file_name, strerror(errno));
		return FALSE;
	}
	app_indicator_set_icon_full(global_app_indicator, temp_file_name, "");
	app_indicator_set_attention_icon_full(global_app_indicator, temp_file_name, "");
	g_bytes_unref(bytes);
	return FALSE;
}

void _systray_menu_item_selected(int *id) {
	systray_menu_item_selected(*id);
}

GtkMenuItem* find_menu_by_id(int id) {
	GList* it;
	for(it = global_menu_items; it != NULL; it = it->next) {
		MenuItemNode* item = (MenuItemNode*)(it->data);
		if(item->menu_id == id) {
			return GTK_MENU_ITEM(item->menu_item);
		}
	}
	return NULL;
}

// runs in main thread, should always return FALSE to prevent gtk to execute it again
gboolean do_add_or_update_menu_item(gpointer data) {
	MenuItemInfo *mii = (MenuItemInfo*)data;
	GList* it;
	for(it = global_menu_items; it != NULL; it = it->next) {
		MenuItemNode* item = (MenuItemNode*)(it->data);
		if(item->menu_id == mii->menu_id) {
			gtk_menu_item_set_label(GTK_MENU_ITEM(item->menu_item), mii->title);

			if (mii->isCheckable) {
				// We need to block the "activate" event, to emulate the same behaviour as in the windows version
				// A Check/Uncheck does change the checkbox, but does not trigger the checkbox menuItem channel
				g_signal_handler_block(GTK_CHECK_MENU_ITEM(item->menu_item), item->signalHandlerId);
				gtk_check_menu_item_set_active(GTK_CHECK_MENU_ITEM(item->menu_item), mii->checked == 1);
				g_signal_handler_unblock(GTK_CHECK_MENU_ITEM(item->menu_item), item->signalHandlerId);
			}
			break;
		}
	}

	// menu id doesn't exist, add new item
	if(it == NULL) {
		GtkWidget *menu_item;
		if (mii->isCheckable) {
			menu_item = gtk_check_menu_item_new_with_label(mii->title);
			gtk_check_menu_item_set_active(GTK_CHECK_MENU_ITEM(menu_item), mii->checked == 1);
		} else {
			menu_item = gtk_menu_item_new_with_label(mii->title);
		}
		int *id = malloc(sizeof(int));
		*id = mii->menu_id;
		long signalHandlerId = g_signal_connect_swapped(
			G_OBJECT(menu_item),
			"activate",
			G_CALLBACK(_systray_menu_item_selected),
			id
		);

		if (mii->parent_menu_id == 0) {
			gtk_menu_shell_append(GTK_MENU_SHELL(global_tray_menu), menu_item);
		} else {
			GtkMenuItem* parentMenuItem = find_menu_by_id(mii->parent_menu_id);
			GtkWidget* parentMenu = gtk_menu_item_get_submenu(parentMenuItem);

			if(parentMenu == NULL) {
				parentMenu = gtk_menu_new();
				gtk_menu_item_set_submenu(parentMenuItem, parentMenu);
			}

			gtk_menu_shell_append(GTK_MENU_SHELL(parentMenu), menu_item);
		}

		MenuItemNode* new_item = malloc(sizeof(MenuItemNode));
		new_item->menu_id = mii->menu_id;
		new_item->signalHandlerId = signalHandlerId;
		new_item->menu_item = menu_item;
		GList* new_node = malloc(sizeof(GList));
		new_node->data = new_item;
		new_node->next = global_menu_items;
		if(global_menu_items != NULL) {
			global_menu_items->prev = new_node;
		}
		global_menu_items = new_node;
		it = new_node;
	}
	GtkWidget* menu_item = GTK_WIDGET(((MenuItemNode*)(it->data))->menu_item);
	gtk_widget_set_sensitive(menu_item, mii->disabled != 1);
	gtk_widget_show(menu_item);

	free(mii->title);
	free(mii->tooltip);
	free(mii);
	return FALSE;
}

gboolean do_add_separator(gpointer data) {
	GtkWidget *separator = gtk_separator_menu_item_new();
	gtk_menu_shell_append(GTK_MENU_SHELL(global_tray_menu), separator);
	gtk_widget_show(separator);
	return FALSE;
}

// runs in main thread, should always return FALSE to prevent gtk to execute it again
gboolean do_hide_menu_item(gpointer data) {
	MenuItemInfo *mii = (MenuItemInfo*)data;
	GList* it;
	for(it = global_menu_items; it != NULL; it = it->next) {
		MenuItemNode* item = (MenuItemNode*)(it->data);
		if(item->menu_id == mii->menu_id){
			gtk_widget_hide(GTK_WIDGET(item->menu_item));
			break;
		}
	}
	return FALSE;
}

// runs in main thread, should always return FALSE to prevent gtk to execute it again
gboolean do_show_menu_item(gpointer data) {
	MenuItemInfo *mii = (MenuItemInfo*)data;
	GList* it;
	for(it = global_menu_items; it != NULL; it = it->next) {
		MenuItemNode* item = (MenuItemNode*)(it->data);
		if(item->menu_id == mii->menu_id){
			gtk_widget_show(GTK_WIDGET(item->menu_item));
			break;
		}
	}
	return FALSE;
}

// runs in main thread, should always return FALSE to prevent gtk to execute it again
gboolean do_quit(gpointer data) {
	_unlink_temp_file();
	// app indicator doesn't provide a way to remove it, hide it as a workaround
	app_indicator_set_status(global_app_indicator, APP_INDICATOR_STATUS_PASSIVE);
	gtk_main_quit();
	return FALSE;
}

void setIcon(const char* iconBytes, int length, bool template) {
	GBytes* bytes = g_bytes_new_static(iconBytes, length);
	g_idle_add(do_set_icon, bytes);
}

void setTitle(char* ctitle) {
	app_indicator_set_label(global_app_indicator, ctitle, "");
	free(ctitle);
}

void setTooltip(char* ctooltip) {
	free(ctooltip);
}

void setMenuItemIcon(const char* iconBytes, int length, int menuId, bool template) {
}

void add_or_update_menu_item(int menu_id, int parent_menu_id, char* title, char* tooltip, short disabled, short checked, short isCheckable) {
	MenuItemInfo *mii = malloc(sizeof(MenuItemInfo));
	mii->menu_id = menu_id;
	mii->parent_menu_id = parent_menu_id;
	mii->title = title;
	mii->tooltip = tooltip;
	mii->disabled = disabled;
	mii->checked = checked;
	mii->isCheckable = isCheckable;
	g_idle_add(do_add_or_update_menu_item, mii);
}

void add_separator(int menu_id) {
	MenuItemInfo *mii = malloc(sizeof(MenuItemInfo));
	mii->menu_id = menu_id;
	g_idle_add(do_add_separator, mii);
}

void hide_menu_item(int menu_id) {
	MenuItemInfo *mii = malloc(sizeof(MenuItemInfo));
	mii->menu_id = menu_id;
	g_idle_add(do_hide_menu_item, mii);
}

void show_menu_item(int menu_id) {
	MenuItemInfo *mii = malloc(sizeof(MenuItemInfo));
	mii->menu_id = menu_id;
	g_idle_add(do_show_menu_item, mii);
}

void quit() {
	g_idle_add(do_quit, NULL);
}
