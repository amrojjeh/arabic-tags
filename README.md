# Project Description
We're working on two projects: NahwApp, and ArabicTagging, which is this one. The purpose of the former is to quiz the student on I'rab, but to do so, it needs
an Arabic corpus that's tagged with grammatical explanations in addition to technical information. This project helps develop that corpus by providing the tools
to tag Arabic text.

There's three phases to this tool:
1) Editing
2) Grammatical tagging
3) Technical tagging

Grammatical and Technical tagging can be done asynchronously. However, editing must be done prior to the last two stages.
As such, the last two stages will be inaccessible until editing is complete, and once editing is complete, it'll be locked.
Unlocking it would require a full reset of both grammatical and technical tagging.

## Editing phase
**Status:** Mostly complete.

The editing phase is where we insert the actual Arabic. The tool ensures that no illegal characters are making it through the corpus,
in addition to preventing double spaces (as well as other forms of whitespace) and short vowels. While it should mostly consist of copy pasting text and minor tweaking,
this tool allows you to also write up your own text very easily within the website.

It's mostly complete, meaning that the central web component is complete, and it just needs to be integrated into the standard workflow.

## Grammatical tagging phase
**Status:** In progress.

This phase is meant to be done by Arabic grammarians. It should let them select individual Arabic words (such as particles and pronouns) and tag them with their traditional
grammatical names (مبتدأ ,خبر, etc...). 

I just started working on this today (08/05/2023), hence why it's "in progress."

## Technical tagging phase
**Status:** Not yet started.

The computer would obviously not be able to understand the traditional grammatical names, which are there to serve as explanations to the student.
So, the technical tagging is meant to tag each word in a way that the computer can understand. For instance, to properly quiz the student on word endings, the program
needs to know which vowels are part of the case endings and which ones can be ignored. The technical tagging phase is meant to provide that information.

## Sharing
**Status:** In progress.

Since it'll likely be rare for a grammarian to revisit this program, we won't be adding any authentication. Instead, anyone is free to create an excerpt, and as long as they don't
share their URL, they will be the only ones allowed to edit their excerpt. If they want to share their excerpt while restricting certain editing privileges, then they can
generate a new URL through the program and set the appropriate constraints. Then, they simply share their "share URL" and Bob's your uncle.

I've already implemented unique and owner URLs, but I've not implemented the share URLS yet. I plan to do so once the workflow is entirely complete and tested.
