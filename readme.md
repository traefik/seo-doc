## SEO handler

This program aims to process a documentation folder from [traefik/doc](https://github.com/traefik/doc) and iterate each HTML file adding the requirements for a better SEO.

### The requirements

1. Older doc versions should contain this meta tag:
```html
<head>
  <!-- ... -->
  <meta name="robots" content="index, nofollow" />
  <!-- ... -->
</head>
```

2. Older doc versions should have a canonical link in the head that points to the latest documentation page. Example:
```html
<!-- in a page under v1.0 -->
<head>
  <!-- ... -->
  <link rel="canonical" href="https://doc.traefik.io/<product_name>" />
  <!-- ... -->
</head>
```

3. Titles in older versions should have the Product name and version as a suffix, and should not have more than 65 characters. For example:
```
Overview | Traefik | v2.0
```

4. sitemap.xml and sitemap.xml.gz should not exist under version folders.

### How to use it

You can use the `seo` directly from command line, and using the path to the documentation dir as parameter.

Examples:

```sh
seo -path /path/to/doc/traefik
seo -path /path/to/doc/traefik-mesh
seo -path /path/to/doc/traefik-pilot
seo -path /path/to/doc/traefik-enterprise
```

```sh
seo -path ./site -product traefik
seo -path ./site -product "traefik-mesh"
seo -path ./site -product "traefik-pilot"
seo -path ./site -product "traefik-enterprise"
```
