const localisation = {
  "en-GB": {
    language: "English",
    username: "Username",
    password: "Password",
    login: "Login",
    credentials_error: "Invalid username or password",
  },
  "es-ES": {
    language: "Español",
    username: "Nombre de usuario",
    password: "Contraseña",
    login: "Iniciar sesión",
    credentials_error: "Nombre de usuario o contraseña incorrectos",
  },
  "fr-FR": {
    language: "Français",
    username: "Nom d'utilisateur",
    password: "Mot de passe",
    login: "Connexion",
    credentials_error: "Nom d'utilisateur ou mot de passe incorrect",
  },
  "ml-IN": {
    language: "മലയാളം",
    username: "ഉപയോക്തൃനാമം",
    password: "രഹസ്യവാക്ക്",
    login: "പ്രവേശിക്കുക",
    credentials_error: "തെറ്റായ ഉപയോക്തൃനാമം അല്ലെങ്കില്‍ രഹസ്യവാക്ക് തെറ്റാണ്",
  },
  "hi-IN": {
    language: "हिन्दी",
    username: "उपयोगकर्ता नाम",
    password: "पासवर्ड",
    login: "लॉग इन करें",
    credentials_error: "गलत उपयोगकर्ता नाम या पासवर्ड",
  },
  "ta-IN": {
    language: "தமிழ்",
    username: "பயனர்பெயர்",
    password: "கடவுச்சொல்",
    login: "உள்நுழைக",
    credentials_error: "தவறான பயனர்பெயர் அல்லது கடவுச்சொல்",
  },
  "bn-IN": {
    language: "বাংলা",
    username: "ব্যবহারকারীর নাম",
    password: "পাসওয়ার্ড",
    login: "লগ ইন",
    credentials_error: "ভুল ব্যবহারকারীর নাম বা পাসওয়ার্ড",
  },
  "zh-CN": {
    language: "简体中文",
    username: "用户名",
    password: "密码",
    login: "登录",
    credentials_error: "用户名或密码错误",
  },
  "zh-TW": {
    language: "繁體中文",
    username: "用戶名",
    password: "密碼",
    login: "登入",
    credentials_error: "用戶名或密碼錯誤",
  },
  "de-DE": {
    language: "Deutsch",
    username: "Benutzername",
    password: "Passwort",
    login: "Anmelden",
    credentials_error: "Ungültiger Benutzername oder Passwort",
  },
  "ru-RU": {
    language: "Русский",
    username: "Имя пользователя",
    password: "Пароль",
    login: "Войти",
    credentials_error: "Неправильное имя пользователя или пароль",
  },
  "pt-BR": {
    language: "Português",
    username: "Nome de usuário",
    password: "Senha",
    login: "Entrar",
    credentials_error: "Nome de usuário ou senha inválidos",
  },
  "ja-JP": {
    language: "日本語",
    username: "ユーザー名",
    password: "パスワード",
    login: "ログイン",
    credentials_error: "ユーザー名またはパスワードが違います",
  },
  "it-IT": {
    language: "Italiano",
    username: "Nome utente",
    password: "Password",
    login: "Accedi",
    credentials_error: "Nome utente o password non validi",
  },
  "tr-TR": {
    language: "Türkçe",
    username: "Kullanıcı adı",
    password: "Şifre",
    login: "Giriş yap",
    credentials_error: "Geçersiz kullanıcı adı veya şifre",
  },
  "ko-KR": {
    language: "한국어",
    username: "사용자 이름",
    password: "비밀번호",
    login: "로그인",
    credentials_error: "잘못된 사용자 이름 또는 비밀번호",
  },
  "vi-VN": {
    language: "Tiếng Việt",
    username: "Tên đăng nhập",
    password: "Mật khẩu",
    login: "Đăng nhập",
    credentials_error: "Tên đăng nhập hoặc mật khẩu không chính xác",
  },
  "id-ID": {
    language: "Bahasa Indonesia",
    username: "Nama pengguna",
    password: "Kata sandi",
    login: "Masuk",
    credentials_error: "Nama pengguna atau kata sandi salah",
  },
};

function localize(language, error = "") {
  const translatedText = localisation[language];
  document.getElementById("username").innerText = translatedText["username"];
  document.getElementById("password").innerText = translatedText["password"];
  document.getElementById("login_button").value = translatedText["login"];
  if (error != "") {
    document.getElementById("login-error").innerText = translatedText[error];
  }
}

// add language options to the select element
const select = document.getElementById("language-select");
for (const [key, value] of Object.entries(localisation)) {
  const option = document.createElement("option");
  option.value = key;
  option.text = value["language"];
  select.appendChild(option);
}
