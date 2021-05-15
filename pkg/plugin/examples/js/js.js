var tagName = "Hawwwwt"

function main() {
    var modeArg = input.Args.mode;
    try {
        if (modeArg == "" || modeArg == "add") {
            console.log("add");
            addTag();
        } else if (modeArg == "remove") {
            console.log("remove");
        } else if (modeArg == "long") {
            console.log("long");
        } else if (modeArg == "indef") {
            console.log("indef");
        }
    } catch (err) {
        return {
            Error: err
        };
    }

    // if err != nil {
    //     errStr := err.Error()
    //     *output = common.PluginOutput{
    //         Error: &errStr,
    //     }
    //     return nil
    // }

    return {
        Output: "ok"
    };
}

function getResult(result) {
    if (result[1]) {
        throw result[1];
    }

    return result[0];
}

function getTagID(create) {
	console.log("Checking if tag exists already")

	// see if tag exists already
    var result = api.AllTags()
    result = getResult(result);

    var tag;
	for (var i = 0; i < result.length; ++i) {
		if (result[i].Name === tagName) {
			tag = result[i];
			break;
		}
	}

    if (tag) {
        console.log("found existing tag");
        return tag.ID;
    }

	if (!create) {
		console.log("Not found and not creating");
		return null;
	}

    console.log("Creating new tag")

	result = api.TagCreate({
        Name: tagName
    });

    result = getResult(result);
    
	return result.ID
}

function addTag() {
    var tagID = getTagID(true)
    
	var scene = findRandomScene();

	if (scene === null) {
		throw "no scenes to add tag to";
	}

	// var m struct {
	// 	SceneUpdate SceneUpdate `graphql:"sceneUpdate(input: $s)"`
	// }

	// input := SceneUpdateInput{
	// 	ID:     scene.ID,
	// 	TagIds: scene.getTagIds(),
	// }

	// input.TagIds = addTagId(input.TagIds, *tagID)

	// vars := map[string]interface{}{
	// 	"s": input,
	// }

	// log.Infof("Adding tag to scene %v", scene.ID)
	// err = client.Mutate(context.Background(), &m, vars)
	// if err != nil {
	// 	return fmt.Errorf("Error mutating scene: %s", err.Error())
	// }

	// return nil
}

function findRandomScene() {
	// get a random scene
    console.log("Finding a random scene")
    result = api.FindScenes(null, [], {
        PerPage: 1,
        Sort: "random"
    });
    result = getResult(result);

    if (result.Count === 0) {
        return null;
    }

    return result.Scenes[0];
}

main();