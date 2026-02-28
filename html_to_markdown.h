#ifndef HTML_TO_MARKDOWN_H
#define HTML_TO_MARKDOWN_H

#ifdef __cplusplus
extern "C" {
#endif

// Convert HTML string to Markdown
char* ConvertHTMLToMarkdown(char* html);

// Convert HTML file content to Markdown
char* ConvertHTMLFileToMarkdown(char* filepath);

#ifdef __cplusplus
}
#endif

#endif // HTML_TO_MARKDOWN_H