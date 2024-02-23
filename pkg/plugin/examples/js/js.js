var tagName = "Hawwwwt"

function main() {
    var modeArg = input.Args.mode;
    if (modeArg !== undefined) {
        try {
            if (modeArg == "" || modeArg == "add") {
                addTag();
            } else if (modeArg == "remove") {
                removeTag();
            } else if (modeArg == "long") {
                doLongTask();
            } else if (modeArg == "indef") {
                doIndefiniteTask();
            } else if (modeArg == "hook") {
                doHookTask();
            }
        } catch (err) {
            return {
                Error: err
            };
        }

        return {
            Output: "ok"
        };
    }

    if (input.Args.error) {
        return {
            Error: input.Args.error
        };
    }

    // immediate mode
    // just return the args
    return {
        Output: input.Args
    };
}

function getResult(result) {
    if (result[1]) {
        throw result[1];
    }

    return result[0];
}

function getTagID(create) {
	log.Info("Checking if tag exists already (via GQL)")

	// see if tag exists already
    var query = "\
query {\
  allTags {\
    id\
    name\
  }\
}"
	
    var result = gql.Do(query);
    var allTags = result["allTags"];
    
    var tag;
	for (var i = 0; i < allTags.length; ++i) {
		if (allTags[i].name === tagName) {
			tag = allTags[i];
			break;
		}
	}

    if (tag) {
        log.Info("found existing tag");
        return tag.id;
    }

	if (!create) {
		log.Info("Not found and not creating");
		return null;
	}

    log.Info("Creating new tag");

    var mutation = "\
mutation tagCreate($input: TagCreateInput!) {\
  tagCreate(input: $input) {\
    id\
  }\
}";

    var variables = {
        input: {
            'name': tagName
        }
    };

    result = gql.Do(mutation, variables);
    log.Info("tag id = " + result.tagCreate.id);
	return result.tagCreate.id;
}

function addTag() {
    var tagID = getTagID(true)
    
	var scene = findRandomScene();

	if (scene === null) {
		throw "no scenes to add tag to";
	}

    var tagIds = []
    var found = false;
	for (var i = 0; i < scene.tags.length; ++i) {
        var sceneTagID = scene.tags[i].id;
        if (tagID === sceneTagID) {
            found = true;
        }
        tagIds.push(sceneTagID);
    }
	
    if (found) {
        log.Info("already has tag");
        return;
    }
	    
    tagIds.push(tagID)

    var mutation = "\
mutation sceneUpdate($input: SceneUpdateInput!) {\
    sceneUpdate(input: $input) {\
        id\
    }\
}";

    var variables = {
        input: {
            id: scene.id,
            tag_ids: tagIds,
        }
    };

    log.Info("Adding tag to scene " + scene.id);

    gql.Do(mutation, variables);
}

function removeTag() {
	var tagID = getTagID(false);

	if (tagID == null) {
		log.Info("Tag does not exist. Nothing to remove");
		return
    }

	log.Info("Destroying tag");
	
    var mutation = "\
mutation tagDestroy($input: TagDestroyInput!) {\
    tagDestroy(input: $input)\
}";

    var variables = {
        input: {
            id: tagID
        }
    };

    gql.Do(mutation, variables);
}

function findRandomScene() {
	// get a random scene
    log.Info("Finding a random scene")

    var query = "\
query findScenes($filter: FindFilterType!) {\
    findScenes(filter: $filter) {\
        count\
        scenes {\
            id\
            tags {\
                id\
            }\
        }\
    }\
}"

    var variables = {
        filter: {
            per_page: 1,
            sort: 'random'
        }
    };

    var result = gql.Do(query, variables);
    var findScenes = result["findScenes"];
    
    if (findScenes.Count === 0) {
        return null;
    }

    return findScenes.scenes[0];
}

function doLongTask() {
	var total = 100;
	var upTo = 0;

	log.Info("Doing long task");
	while (upTo < total) {
		util.Sleep(1000);

		log.Progress(upTo / total);
		upTo = upTo + 1;
    }
}

function doIndefiniteTask() {
	log.Info("Sleeping indefinitely");
	while (true) {
		util.Sleep(1000);
    }
}

function doHookTask() {
    log.Info("JS Hook called!");
    log.Info(input.Args);
}

main();