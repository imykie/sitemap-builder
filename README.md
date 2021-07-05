# Sitemap Builder

This is a tool that maps all the pages within a specific website domain depending on the specified depth.

Read more about sitemaps [here](https://www.sitemaps.org/protocol.html)
### Flags

  |S/N | Flag  | Type  | Required/Optional   |  Description | Default|
  |---|---|---|---|---|---|
  | 1  | url | String | Optional  | The URL you want to build Sitemap for  | https://go.dev |
  |  2 | depth  | Int  | Optional  | The maximum depth of the Sitemap Builder  | 3 |

### Usage
Ensure you are at the project root:

##### Printing generated sitemap xml to console.

```shell
 go build . && go run . --url="example.com" --depth=3
```

##### Creating a xml file for the generated sitemap xml
```shell
 go build . && go run . --url="example.com" --depth=3 > sitemap.xml
```

### Result

```xml
<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url>
    <loc>http://www.example.com/</loc>
  </url>
  <url>
    <loc>http://www.example.com/sample</loc>
  </url>
  <!--Total number of URLs: 2-->
</urlset>
```