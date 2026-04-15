# Tutorial: Integrating Text File Support into Stash

## Introduction

In this comprehensive tutorial, you'll learn how to add text file support to Stash, allowing you to manage and organize your text files with the same level of detail as image or video content. The guide will walk you through setting up the necessary metadata fields, configuring BBCode for text formatting, linking audio files to text, and playing both simultaneously.

### What You Will Learn

- How to add support for text files in Stash.
- Implementing metadata fields such as name, author, URL, and more.
- Configuring BBCode for formatted text.
- Linking and synchronizing audio with text content.

## Prerequisites

Before proceeding, ensure you have the following:

1. **Stash Application**: Ensure that you are using the latest version of Stash available on GitHub or your preferred package manager.
2. **Development Environment**: A local development environment set up to work with Git and any necessary IDEs (e.g., Visual Studio Code).
3. **Basic Understanding of JavaScript and CSS**: Familiarity with these technologies will be beneficial as you integrate the new features.

## Step 1: Setting Up Metadata Fields

### Adding Metadata to Text Files

To enable metadata support, we need to add custom fields to our text files. We'll use a JSON-like structure for simplicity. Here's an example of how your text file might look:

```json
{
    "name": "Example Story",
    "author": "John Doe",
    "url_for_source": "https://example.com/source",
    "website_of_source": "https://example.com/website",
    "date": "2023-10-01",
    "language": "English",
    "cover_front": "/path/to/cover_front.jpg",
    "cover_back": "/path/to/cover_back.jpg",
    "rating": 4,
    "tags": ["fiction", "adventure"],
    "details": "A thrilling adventure story.",
    "tag_line": "Explore the unknown.",
    "link_to_audio": "/path/to/audio.mp3"
}
```

### Creating Metadata Files

Ensure that each text file has a corresponding metadata file with the same name and a `.json` extension. For example, if your text file is named `example_story.txt`, you should have an accompanying `example_story.json`.

## Step 2: Integrating BBCode Support

### Enabling BBCode Parsing

To enable BBCode support for text formatting, we need to add a custom filter in Stash's configuration.

1. **Locate the Configuration File**:
   - Navigate to your Stash installation directory and locate the `config.js` file.
2. **Edit the Config File**:

```javascript
// config.js

const stash = {
    // existing configurations...
    
    textFileSupport: {
        bbCodeEnabled: true,
        supportedTags: [
            'b', 'i', 'u', 's', 'sub', 'sup', 'quote', 'code'
        ]
    }
};
```

3. **Save and Restart**:
   - Save the changes and restart Stash to apply the new configuration.

## Step 3: Linking Audio Files

### Synchronizing Text with Audio

To link an audio file to a text file, you need to ensure that both files are correctly named and placed in the appropriate directory. The `link_to_audio` field in your JSON metadata should point to the correct path of the audio file.

1. **Place Audio File**:
   - Ensure that the audio file is stored in the same directory as its corresponding text file.
2. **Update Metadata**:
   - Make sure the `link_to_audio` field points to the correct path, e.g., `/path/to/audio.mp3`.

## Step 4: Playing and Reading Text with Audio

### Implementing Play Functionality

To play both the text and audio simultaneously, you need to integrate a player component that can read from the text file and play the corresponding audio.

1. **Create a Player Component**:
   - You may use existing libraries like `react-audio-player` or create your own using HTML5 `<audio>` elements.
2. **Integrate with Stash**:

```javascript
// In your JavaScript component

import React, { useState, useEffect } from 'react';
import AudioPlayer from 'react-h5-audio-player';

function TextAudioPlayer({ filePath }) {
    const [text, setText] = useState('');
    const [audioSrc, setAudioSrc] = useState('');

    useEffect(() => {
        // Load text and audio paths
        fetchMetadata(filePath);
    }, [filePath]);

    const fetchMetadata = async (filePath) => {
        try {
            const metadataPath = filePath.replace('.txt', '.json');
            const response = await fetch(metadataPath);
            const data = await response.json();

            if (data.link_to_audio) {
                setAudioSrc(data.link_to_audio);
            }

            setText('');
            // Implement logic to read and display text from the file
        } catch (error) {
            console.error('Error fetching metadata:', error);
        }
    };

    return (
        <div>
            <h1>{text}</h1>
            <AudioPlayer src={audioSrc} />
        </div>
    );
}
```

3. **Display in Stash**:
   - Integrate the `TextAudioPlayer` component into your Stash UI where text and audio files are displayed.

## Troubleshooting

### Issue: Metadata Not Displaying Correctly
- Ensure that the metadata file is correctly named and located next to the text file.
- Check for any syntax errors in your JSON metadata or JavaScript code.

### Issue: Audio Player Not Playing
- Verify that the path provided in `link_to_audio` is correct.
- Ensure that the audio file format is supported by Stash's current setup.

## Conclusion

By following this comprehensive tutorial, you have successfully added support for text files and linked them with audio content. This integration enhances the functionality of Stash, making it a more versatile tool for managing multimedia assets.

If you encounter any issues or need further assistance, please consult the official Stash documentation or seek help from the community forums.