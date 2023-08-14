from time import time

now = time()

import json
from camel_tools.utils.normalize import normalize_unicode
from camel_tools.utils.dediac import dediac_ar
from camel_tools.tokenizers.word import simple_word_tokenize
from camel_tools.disambig.bert import BERTUnfactoredDisambiguator

bert = BERTUnfactoredDisambiguator.pretrained(top=10)

# Average delta = ~14 seconds
print(f"Ready! {time() - now:.2f} seconds")


while True:
	words = []
	sentence = input()
	if sentence == "exit":
		break
	sentence = simple_word_tokenize(dediac_ar(normalize_unicode(sentence)))
	disambig = bert.disambiguate(sentence)

	for w in disambig:
		for sa in w.analyses:
			print((sa.score, sa.analysis["diac"], sa.analysis["cas"]))
		print()
