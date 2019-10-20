#!/usr/bin/env python3
# -*- coding: utf-8 -*-
# Find out magnet links inside pages
import re
import config, event

class Analyzer():
	'Analyzer(): Extract magnet links and title from page content.'

	def get_magnet_links(self, result_text):
		'Get magnet links from page content. Return a list containing strings of magnet links.\n\nresult_text: Page content in utf-8 encoding.\n\nEvent:\n\tAnalyzingMangetLinks(magnets)'
		if config.analyzer['search_content_only']:
			match = re.search(config.analyzer['regex_content'], result_text, re.MULTILINE | re.DOTALL)
			result_text = match.group(1) if match else ''

		hashes = list(set(re.findall(config.analyzer['regex_magnets_40'], result_text))) # 40-char magnets
		hashes.extend(list(set(re.findall(config.analyzer['regex_magnets_32'], result_text)))) # 32-char magnets
		magnets = list(set([(config.analyzer['magnet_prefix'] + hash_value).lower() for hash_value in hashes]))

		event.event.dispatch('AnalyzingMangetLinks', magnets = magnets)
		return magnets

	def get_page_title(self, result_text):
		'Get page title from page content. Return title string or empty string if it is not found.\n\nresult_text: Page content in utf-8 encoding.\n\nEvent:\n\tAnalyzingTitle(title)'
		match = re.search(config.analyzer['regex_title'], result_text)
		if match:
			title = match.group(1).strip()
		else:
			title = ''
		event.event.dispatch('AnalyzingTitle', title = title)
		return title
