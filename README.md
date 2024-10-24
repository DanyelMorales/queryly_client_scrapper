
# Queryly CLI Tool

## Overview

This CLI tool fetches news articles in bulk from different journal sources, providing various options to customize the fetching process. It supports features such as fetching articles by portal, displaying available portals, and retrieving news in different formats and configurations.

## Features

- **Fetch articles**: Scrape and download articles from various online journal portals.
- **Display available IDs**: List the portal IDs available for scraping.
- **Display options**: View scraping options for specific portals.
- **Cron support**: Manage periodic scraping tasks. (it's available but we are still documenting it)
- **Health check**: Perform a health check on the tool. (it's available but we are still documenting it)
- **Version command**: Check the current version of the scraper.

## Manual Installation

1. Clone the repository and navigate to the project folder.
2. Install the necessary dependencies using Go.

```bash
go mod tidy
go build
```

3. Add the binary to your `$PATH` or create a symbolic link for easy access:

```bash
ln -s /path/to/binary /usr/local/bin/querylyctl
```

### Makefile installation

### `build`
Runs the entire build process, including setting up directories, compiling the code, compressing the binaries, and cleaning up temporary files.

```bash
make build
```

**Process:**
1. **Setup**: Creates the necessary directories (`build` and `bin`).
2. **Compile**: Downloads dependencies, sets version information, compiles the Go source code, and prepares the binary.
3. **Compress**: Packages the binary and dependencies into a `.tar.xz` file, generates an installer, and compresses it into a `.zip` file.
4. **Clean**: Removes temporary build files after the process is completed.

### `build-run`
Combines the `build` and `install` steps, then runs the CLI tool's version command to verify that the tool was successfully installed.

```bash
make build-run
```

**Process:**
1. Executes the `build` target to compile and package the tool.
2. Installs the tool by running the `setup.sh` script.
3. Runs the tool to display the current version.

## Usage

### Root Command

```bash
querylyctl [command] [flags]
```

The root command provides the main interface to the tool, allowing you to fetch articles, display portal IDs, and access other features.

### Commands

#### `fetch`
Fetch articles from a specified portal.

```bash
querylyctl articles fetch --portal-id [PORTAL_ID]  
```

**Flags**:
- `--id`: The ID of the portal to fetch articles from (required).
- `--end-index`: The page number to navigate (default is `1`).
- `--batch-size`: The number of articles to fetch (default is `10`).
- `--section`: Fetch articles from a specific section.
- `--query`: Search for a specific term.
- `--out`: Save results inside a new subdirectory.
- `--date`: sort the search by Date, values: (last-week, today, last-month).

#### `display-ids`
Display available portal IDs for fetching.

```bash
querylyctl  --cfg ./cfg.json    articles --ids
```

#### `display-options`
Display available scraping options for a specific portal.

```bash
querylyctl articles --cfg ./cfg.json   options --id [PORTAL_ID]
```

**Flags**:
- `--id`: The ID of the portal to display options for (required).

#### `cron`
Manage scheduled scraping tasks (requires additional setup).

#### `health-check`
Run a health check to verify the tool's functionality.

```bash
querylyctl health-check
```

#### `version`
Display the current version of the tool.

```bash
querylyctl version
```

## Configuration

The tool uses a configuration file for some operations. You can specify the file path using the `--config` flag:

```bash
querylyctl --cfg ./cfg.json 
```

If no configuration file is provided, the tool will prompt you to enter settings.

## Example

Fetch articles from all portals using the following command:

```bash
querylyctl articles --cfg ./cfg.json  fetch --id="all" --end-index 5 ---section "Technology"
```

Scrape articles from all portals, fetching only 4 articles for each portal, matching the keyword in the query flag
with a creation date of today.

```bash
newsctl articles --cfg ./cfg.json --id=all  fetch --batch-size=4 --query="keysearch" --date=today
```
Scrape articles from all portals, fetching only 4 articles for each portal, matching the keyword in the query flag
with a creation date of last week.

```bash
newsctl articles --cfg ./cfg.json --id=all  fetch --batch-size=4 --query="keysearch" --date=last-week
```
Scrape articles from all portals, fetching only 4 articles for each portal, matching the keyword in the query flag
with a creation date of last month.

```bash
newsctl articles --cfg ./cfg.json --id=all  fetch --batch-size=4 --query="keysearch" --date=last-month
```
Scrape articles from website1,website2,website3 portals; fetching only 4 articles for each portal, and matching the keyword in the query flag
with a creation date of last month.
```bash
newsctl articles --cfg ./cfg.json  --id="website1,website2,website3"  fetch --batch-size=4 --query="keyword" --date=last-week --out="./"
```
## Configuration file

This JSON configuration file defines settings for the cli tool that is designed to collect and store information from specified websites. It includes the operating mode, paths for storing data, and registry configurations for different websites to scrape.

### Structure

The configuration is organized as follows:

- `roboto`: The main object that houses all the configurations for the scraper.
    - `mode`: Defines the operating mode of the scraper.
    - `scrapperConfig`: Contains settings related to scraping functionality, including output paths.
    - `registry`: A list of sites with individual configurations for scraping.

---
 

##### `outputPath` (String)
- **Description**: Specifies the path where the scraped data will be stored.
- **Default**: `"~/collected/"`
- **Example**:
  ```json
  "outputPath": "~/collected/"
  ```

---

### `registry` (Array of Objects)
- **Description**: A list of websites to scrape, with configurations for each site.

Each object in the registry array defines the parameters required to scrape a specific website.

#### `Context` (String)
- **Description**: A user-defined label representing the name or context of the website to scrape.
- **Example**:
  ```json
  "Context": "NameOfTheSide"
  ```

#### `Enabled` (Boolean)
- **Description**: Determines whether the scraper is active for this specific website.
- **Values**:
    - `true`: Scraping is enabled.
    - `false`: Scraping is disabled.
- **Example**:
  ```json
  "Enabled": true
  ```

#### `Type` (Integer)
- **Description**: Indicates the type of scraper to use for this website. This could represent different scraping methods or strategies, based on internal definitions.
- **Allowed Values**:
    - `0`: Queryly
- **Example**:
  ```json
  "Type": 0
  ```

#### `Host` (String)
- **Description**: The base URL of the website to scrape.
- **Example**:
  ```json
  "Host": "https://website_to_scrape.com/"
  ```

#### `ApiKey` (String)
- **Description**: The API key used for authenticating queries to the website. This is used when the site provides an API endpoint.
- **Example**:
  ```json
  "ApiKey": "queryly_API_KEY"
  ```

#### `Selector` (String)
- **Description**: A CSS selector used to identify specific HTML elements on the webpage to extract the relevant content (e.g., articles).
- **Example**:
  ```json
  "Selector": ".article-custom-selector"
  ```

#### `OverrideExistingNews` (Boolean)
- **Description**: Determines whether the scraper should overwrite existing news articles that are already stored.
- **Allowed Values**:
    - `true`: The scraper will overwrite previously saved news articles.
    - `false`: The scraper will skip already existing news articles.
- **Example**:
  ```json
  "OverrideExistingNews": true
  ```

---

### Example Configuration

```json
{
  "roboto": {
    "mode": "master",
    "scrapperConfig": {
      "outputPath": "~/collected/"
    },
    "registry": [
      {
        "Context": "NameOfTheSide",
        "Enabled": true,
        "Type": 0,
        "Host": "https://website_to_scrape.com/",
        "ApiKey": "queryly_API_KEY",
        "Selector": ".article-custom-selector",
        "OverrideExistingNews": true
      }
    ]
  }
}
```

 ## Contributing

Contributions are welcome! Feel free to submit a pull request or open an issue for bug reports, feature requests, or improvements.

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.