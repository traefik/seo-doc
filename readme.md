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
  <link rel="canonical" href="https://doc.traefik.io/" />
  <!-- ... -->
</head>
```

3. Titles in older versions should have the Product name and version as a suffix, and should not have more than 65 characters. For example:
```
  Overview | Traefik | v2.0
```

4. Sitemap.xml and Sitemap.xml.gz should not exist under version folders.

5. Latest documentation pages (not older) should have a meta description with a brief summary of what that page is about. If the md file has an input hidden with ID 'meta-description', it will be automatically promoted as a meta tag. It's recommended that it has less than 156 characters.

How it should be in the documentation md file:
```html
<input type="hidden" id="meta-description" value="This article explains how to configure a router...">
```

How it should be in the final HTML file:
```html
<head>
  <!-- ... -->
  <meta name="description" content="This article explains how to configure a router..." />
  <!-- ... -->
</head>
```

### How to use it

You can use the `seo` directly from command line, and using the path to the documentation dir as parameter. Example:

```sh
seo /path/to/doc/traefik
seo /path/to/doc/traefik-mesh
seo /path/to/doc/traefik-pilot
seo /path/to/doc/traefik-enterprise
```
