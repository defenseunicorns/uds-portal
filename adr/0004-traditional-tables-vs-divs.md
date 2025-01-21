# 4 Traditional Table Structure vs Divs

## Status

Accepted (3 Dec 2024)

## Context

Traditional HTML tables were being used to display our data. For most tables this worked well as tables are made for tabular data; however, in our case we are rendering Svelte components in the table cells, which caused the following issues with responsiveness:

- Table cells that contained components did not resize consistently
- Cutting off text became a problem with certain types of data and components
- We had no control over how to display certain data inside of the cells because the browser determined the width of each table cell based on with width of the content

We need the ability to manipulate the layout more and that is the key. Tables were not made for layout purposes. Using a `div` or any other block element (i.e. `section` etc) allows us to have more flexibility to decide how we want to display the table cells. We have the ability to set widths to table cells and decide how we want to display the data in that given table cell. We can choose to wrap long text or cut off the text using ellipsis. This flexibility allows us to have more options when changing screen sizes. As screen size decreases, we need more flexibility to make the content behave the way we want.

## Decision

Moving forward we will be using `div`'s or block elements that are not the traditional HTML `<table>` for displaying the data.

## Rationale

### Semantics

Tables main purpose are for displaying tabular data. The original spec [here](https://www.w3.org/TR/2018/SPSD-html32-20180315/#table) says:

> HTML 3.2 includes a widely deployed subset of the specification given in RFC 1942 and can be used to markup tabular material or for layout purposes. Note that the latter role typically causes problems when rending to speech or to text only user agents.

Accessibility is an up coming area that we will be addressing and so not only is it semantically incorrect to use tables for layouts, tables are not great for accessibility best practices and future proofing a website.

Using semantic tags on your website and leads to better SEO performance which is another topic we will be addressing soon.

### Customization of Layouts

Tables can be used for layouts as mentioned in the spec; however, tables can break quickly when trying to customize layouts and controlling how content inside a table cell works. Many different solutions have been provided for us to use with other HTML elements. CSS solutions such as flex-box and grid that allow us to manipulate the layout of web pages.

The reason we needed to switch to using a table-less layout no our site, even though most of our data is tabular, is because of the need to be more responsive as well as having custom components that needed to be housed inside of the table cells. We needed the ability to control how those components are displayed and we also needed the ability to easily re-size content and display truncation in the form of ellipsis or using overflow to allow for horizontal scrolling inside a table row or table cell.

With a table-less layout we are allowed to set the size of each column and control how the content inside the cells of those columns behave.

## Consequences

### Positive:

- Control over layouts
- Control over content inside of "table cells"
- Better accessibility
- Better SEO

### Negative:

- More maintenance due to having to set the size of each columns
- Requires more advanced knowledge of CSS/ HTML
