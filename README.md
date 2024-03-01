# Project Description
We're working on two projects: NahwApp, and ArabicTagging, which is this one. The purpose of the former is to quiz the student on I'rab, but to do so, it needs
an Arabic corpus that's tagged with grammatical explanations in addition to technical information. This project helps develop that corpus by providing the tools
to tag Arabic text.

There are two phases to this tool, which must be done in order:
1) Manuscript
2) Editing

## Manuscript phase
The editing phase is where we insert the actual Arabic. The tool ensures that no illegal characters are making it through the corpus,
in addition to preventing double spaces (as well as other forms of whitespace) and short vowels. While it should mostly consist of copy pasting text and minor tweaking,
this tool allows you to also write up your own text very easily within the website.

## Editing phase
The text from the Manuscript phase will be given to an AI to vowelize and tokenize (credits to CAMeL Lab). After that, the user can correct the short vowels and insert
grammatical data for each word. Once the text is ready, it could be exported as a JSON and be used in any project. At the moment, however, only NahwApp can read the data.

