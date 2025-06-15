# Anytype Readwise Sync

This script synchronizes your bookmarks from Readwise/Reader to your Anytype space. It fetches all your Readwise bookmarks and creates a new object for each one in Anytype, using a specified template for the content.

## Usage

The main purpose of this script is to automate the process of saving your Readwise bookmarks into your Anytype knowledge base. You can customize the type of object created, the space it's saved in, and the template used for its content.



## Setup and Configuration

### 1. Environment Variables

Before running the script, you need to configure your credentials. Create a `.env` file in the root directory of the project with the following content:
```
READWISE_TOKEN=your_readwise_access_token
ANYTYPE_API_KEY=your_anytype_api_key
```
**Optional:**
```
ANYTYPE_API_BASE_URL=http://localhost:31009
ANYTYPE_VERSION=2025-05-20
```



-   `READWISE_TOKEN`: Your Readwise Access Token. You can get one [here](https://readwise.io/access_token).
-   `ANYTYPE_API_KEY`: Your Anytype API Key.
-   `ANYTYPE_API_BASE_URL` (Optional): The base URL for the Anytype API. Defaults to `http://localhost:31009`.
-   `ANYTYPE_VERSION` (Optional): The Anytype API version. Defaults to `2025-05-20`.

### 2. Templates

You can define the structure of the Anytype objects using one of two methods:

#### Markdown Template
Create a local Markdown file (e.g., `book_template.md`). The script will use the content of this file as the template. You must include the `%%Content%%` placeholder, which will be replaced with the Readwise bookmark's content.

#### Anytype Template (Not yet implemented)
Alternatively, you can use an existing template from your Anytype space. You will need the ID of the template object.

## How to Run

Execute the script from your terminal. You can use flags to configure its behavior.

### Command-line Flags

-   `-template`: Path to the markdown template file (default: `book_template.md`).
-   `-anytype-template`: The ID of an Anytype template object. If provided, it overrides the local markdown template.
-   `-type`: The type of Anytype object to create (default: `Bookmark`).
-   `-space`: The ID of the Anytype space where objects will be created. (default: First space in the list)

### Examples

**Basic Run (using a markdown template):**
This will create objects of type `Bookmark` using `book_template.md` as the template.

*Using a specific markdown template and object type:*
```bash 
go run main.go -template="my_article_template.md" -type="Article"
``` 

**Using an Anytype Template and specifying a Space:** This will use a specific space

```bash
go run main.go -space="<your-space-id>"
```


## Limitations

-   **Sync State**: The script uses the `description` property of the created Anytype object to store a unique identifier for the Readwise bookmark. **Do not modify or remove the content of the `description` field** in the generated objects. If you do, the script will lose track of the synced item and create a duplicate on the next run.
- **Already Sync**: If a object has already been sync the AnyType API doesn't allow updating its body.
- **Anytype Templates**: The integration with the AnyType template system hasn't been done yet, but it can be promising.
- **Cover Image** Currently, there's no way to set a background cover for the article's image