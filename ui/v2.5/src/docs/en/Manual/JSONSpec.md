# Import/Export JSON Specification

The metadata given to Stash can be exported into the JSON format. This structure can be modified, or replicated by other means. The resulting data can then be imported again, giving the possibility for automatic scraping of all kinds. The format of this metadata bulk is a folder structure, containing the following folders:
  
* `files`
* `galleries`
* `images`
* `performers`
* `scenes`
* `studios`
* `groups`

## File naming

When exported, files are named with different formats depending on the object type:

| Type | Format |
|------|--------|
| Files/Folders | `<path depth in hex, two character width>.<basename>.<hash>.json` |
| Galleries | `<first zip filename>.<path hash>.json` or `<folder basename>.<path hash>.json` or `<title>.json` |
| Images | `<title or first file basename>.<hash>.json` |
| Performers | `<name>.json` |
| Scenes | `<title or first file basename>.<hash>.json` |
| Studios | `<name>.json` |
| Groups | `<name>.json` |

Note that the file naming is not significant when importing. All json files will be read from the subdirectories.
  
## Content of the json files

In the following, the values of the according jsons will be shown. If the value should be a number, it is written with after comma values (like `29.98` or `50.0`), but still as a string. The meaning from most of them should be obvious due to the previous explanation or from the possible values stash offers when editing, otherwise a short comment will be added.

The json values are given as strings, if not stated otherwise. Every new line will stand for a new value in the json. If the value is a list of objects, the values of that object will be shown indented.  

If a value is empty in any file, it can be left out of the file entirely. 
Many files have an `created_at` and `updated_at`, both are kept in the following format:
```  
YYYY-MM-DDThh:mm:ssTZD  
```
Example:  
```
"created_at": "2019-05-03T21:36:58+01:00"
```

### Performer
```
name  
url  
twitter  
instagram  
birthdate  
death_date  
ethnicity  
country  
hair_color  
eye_color  
height  
weight  
measurements  
fake_tits  
career_length  
tattoos  
piercings  
image (base64 encoding of the image file)  
created_at  
updated_at
rating (integer)
details
```

### Studio
```
name  
url  
image (base64 encoding of the image file)  
created_at  
updated_at
rating (integer)  
details  
```

### Scene
```
title  
studio  
url  
date  
rating (integer)  
details  
performers (list of strings, performers name)  
tags (list of strings)  
markers     
  title  
  seconds  
  primary_tag  
  tags (list of strings)  
  created_at  
  updated_at  
file (not a list, but a single object)  
  size (in bytes, no after comma values)  
  duration (in seconds)  
  video_codec (example value: h264)  
  audio_codec (example value: aac)  
  width (integer, in pixel)  
  height (integer, in pixel)  
  framerate  
  bitrate (integer, in Bit)  
created_at  
updated_at  
```


### Image
```
title  
studio  
rating (integer)  
performers (list of strings, performers name)  
tags (list of strings)  
files (list of path strings)
galleries
  zip_files (list of path strings)
  folder_path
  title (for user-created gallery)
created_at  
updated_at  
```

### Gallery
```
title  
studio  
url  
date  
rating (integer)  
details  
performers (list of strings, performers name)  
tags (list of strings)  
zip_files (list of path strings)
folder_path   
created_at  
updated_at  
```

## Files

### Folder
```
zip_file (path to containing zip file)
mod_time
type (= folder)
path
created_at
updated_at
```

### Video file
```
zip_file (path to containing zip file)
mod_time
type (= video)
path
fingerprints
  type
  fingerprint
size
format
width
height
duration
video_codec
audio_codec
frame
bitrate
interactive (bool)
interactive_speed (integer)
created_at
updated_at
```

### Image file
```
zip_file (path to containing zip file)
mod_time
type (= image)
path
fingerprints
  type
  fingerprint
size
format
width
height
created_at
updated_at
```

### Other files
```
zip_file (path to containing zip file)
mod_time
type (= file)
path
fingerprints
  type
  fingerprint
size
created_at
updated_at
```

## In JSON format

For those preferring the json-format, defined [here](https://json-schema.org/), the following format may be more interesting:

### performer.json

``` json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "$id": "https://docs.stashapp.cc/in-app-manual/tasks/jsonspec#performerjson",
  "title": "performer",
  "description": "A json file representing a performer. The file is named by a MD5 Code.",
  "type": "object",
  "properties": {
    "name": {
      "description": "Name of the performer",
      "type": "string"
    },
    "url": {
      "description": "URL to website of the performer",
      "type": "string"
    },
    "twitter": {
      "description": "Twitter name of the performer",
      "type": "string"
    },
    "instagram": {
      "description": "Instagram name of the performer",
      "type": "string"
    },
    "birthdate": {
      "description": "Birthdate of the performer. Format is YYYY-MM-DD",
      "type": "string"
    },
    "death_date": {
      "description": "Death date of the performer. Format is YYYY-MM-DD",
      "type": "string"
    },
    "ethnicity": {
      "description": "Ethnicity of the Performer. Possible values are black, white, asian or hispanic",
      "type": "string"
    },
    "country": {
      "description": "Country of the performer",
      "type": "string"
    },
    "hair_color": {
      "description": "Hair color of the performer",
      "type": "string"
    },
    "eye_color": {
      "description": "Eye color of the performer",
      "type": "string"
    },
    "height": {
      "description": "Height of the performer in centimeters",
      "type": "string"
    },
    "weight": {
      "description": "Weight of the performer in kilograms",
      "type": "string"
    },
    "measurements": {
      "description": "Measurements of the performer",
      "type": "string"
    },
    "fake_tits": {
      "description": "Whether performer has fake tits. Possible are Yes or No",
      "type": "string"
    },
    "career_length": {
      "description": "The time the performer has been in business. In the format YYYY-YYYY",
      "type": "string"
    },
    "tattoos": {
      "description": "Giving a description of Tattoos of the performer if any",
      "type": "string"
    },
    "piercings": {
      "description": "Giving a description of Piercings of the performer if any",
      "type": "string"
    },
    "image": {
      "description": "Image of the performer, parsed into base64",
      "type": "string"
    },
    "created_at": {
      "description": "The time this performers data was added to the database. Format is YYYY-MM-DDThh:mm:ssTZD",
      "type": "string"
    },
    "updated_at": {
      "description": "The time this performers data was last changed in the database. Format is YYYY-MM-DDThh:mm:ssTZD",
      "type": "string"
    },
    "details": {
      "description": "Description of the performer",
      "type": "string"
    }
  },
  "required": ["name", "ethnicity", "image", "created_at", "updated_at"]
}

```

### studio.json

``` json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "$id": "https://docs.stashapp.cc/in-app-manual/tasks/jsonspec#studiojson",
  "title": "studio",
  "description": "A json file representing a studio. The file is named by a MD5 Code.",
  "type": "object",
  "properties": {
    "name": {
      "description": "Name of the studio",
      "type": "string"
    },
    "url": {
      "description": "URL to the studios websites",
      "type": "string"
    },
    "image": {
      "description": "Logo of the studio, parsed into base64",
      "type": "string"
    },
    "created_at": {
      "description": "The time this studios data was added to the database. Format is YYYY-MM-DDThh:mm:ssTZD",
      "type": "string"
    },
    "updated_at": {
      "description": "The time this studios data was last changed in the database. Format is YYYY-MM-DDThh:mm:ssTZD",
      "type": "string"
    },
    "details": {
      "description": "Description of the studio",
      "type": "string"
    }
  },
  "required": ["name", "image", "created_at", "updated_at"]
}
```

### scene.json

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "$id": "https://docs.stashapp.cc/in-app-manual/tasks/jsonspec#scenejson",
  "title": "scene",
  "description": "A json file representing a scene. The file is named by the MD5 Code of the file its data is referring to.",
  "type": "object",
  "properties": {
    "title": {
      "description": "Title of the scene",
      "type": "string"
    },
    "studio": {
      "description": "The name of the studio that produced that scene",
      "type": "string"
    },
    "url": {
      "description": "The url to the scenes original source",
      "type": "string"
    },
    "date": {
      "description": "The release date of the scene. Its given in the format YYYY-MM-DD",
      "type": "string"
    },
    "rating": {
      "description": "The scenes Rating. Its given in stars, from 1 to 5",
      "type": "integer"
    },
    "details": {
      "description": "A description of the scene, containing things like the story arc",
      "type": "string"
    },
    "performers": {
      "description": "A list of names of the performers in this gallery",
      "type": "array",
      "items": {
        "type": "string"
      },
      "minItems": 1,
      "uniqueItems": true
    },
    "tags": {
      "description": "A list of the tags associated with this scene",
      "type": "array",
      "items": {
        "type": "string"
      },
      "minItems": 1,
      "uniqueItems": true
    },
    "markers": {
      "description": "Markers mark certain events in the scene, most often the change of the position. They are attributed with their own tags.",
      "type": "array",
      "items": {
        "type": "object",
        "properties": {
          "title": {
            "description": "Searchable name of the marker",
            "type": "string"
          },
          "seconds": {
            "description": "At what second the marker is set. It is given with after comma values, such as 10.0 or 17.5",
            "type": "string"
          },
          "primary_tag": {
            "description": "A tag identifying this marker. Multiple markers from the same scene with the same primary tag are concatenated, showing them as similar in nature",
            "type": "string"
          },
          "tags": {
            "description": "A list of the tags associated with this marker",
            "type": "array",
            "items": {
              "type": "string"
            },
            "minItems": 1,
            "uniqueItems": true
          },
          "created_at": {
            "description": "The time this marker was added to the database. Format is YYYY-MM-DDThh:mm:ssTZD",
            "type": "string"
          },
          "updated_at": {
            "description": "The time this marker was updated the last time. Format is YYYY-MM-DDThh:mm:ssTZD",
            "type": "string"
          }

        },
        "required": ["seconds", "primary_tag", "created_at", "updated_at"]
      },
      "minItems": 1,
      "uniqueItems": true
    },
    "files": {
      "description": "A list of paths of the files for this scene",
      "type": "array",
      "items": {
        "type": "string"
      },
      "minItems": 1,
      "uniqueItems": true
    },
    "created_at": {
      "description": "The time this studios data was added to the database. Format is YYYY-MM-DDThh:mm:ssTZD",
      "type": "string"
    },
    "updated_at": {
      "description": "The time this studios data was last changed in the database. Format is YYYY-MM-DDThh:mm:ssTZD",
      "type": "string"
    }
  },
  "required": ["files", "created_at", "updated_at"]
}
```
