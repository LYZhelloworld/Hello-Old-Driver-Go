#!/usr/bin/env python3
# -*- coding: utf-8 -*-
import sys, os.path
import scanner, analyzer, output
import event

# Event handlers
def event_scanning_page(**kwargs):
	print('Entering: {0}'.format(kwargs['url']))

def event_magnets(**kwargs):
	print('\tFound {0} magnet link(s).'.format(len(kwargs['magnets'])))

def event_removing_duplicates(**kwargs):
	print('Removing duplicated resources...')

def event_before_saving_json(**kwargs):
	print('Saving to JSON file...', end = '')

def event_after_saving(**kwargs):
	print('Done.')

def event_before_saving_txt(**kwargs):
	print('Saving to TXT file...', end = '')

if __name__ == '__main__':
	s = scanner.Scanner(domain = 'www.llss.pw')
	a = analyzer.Analyzer()
	o = output.Output()
	# Read last progress
	if os.path.isfile('result.json'):
		with open('result.json', 'r', encoding = 'utf8') as f:
			o.load_json(f)
			s.load_visited_urls(o.resource_list)
	# Subscribe events
	event.event.subscribe('ScanningPage', event_scanning_page)
	event.event.subscribe('AnalyzingMangetLinks', event_magnets)
	event.event.subscribe('RemovingDuplicatedResources', event_removing_duplicates)
	event.event.subscribe('BeforeSavingJSON', event_before_saving_json)
	event.event.subscribe('BeforeSavingTXT', event_before_saving_txt)
	event.event.subscribe('AfterSavingJSON', event_after_saving)
	event.event.subscribe('AfterSavingTXT', event_after_saving)
	# Loop
	try:
		for url, result_text in s:
			error = False
			if not result_text:
				# Error page
				magnets = []
				title = ''
				error = True
			else:
				magnets = a.get_magnet_links(result_text)
				title = a.get_page_title(result_text)
			o.add(title = title, url = url, magnets = magnets, error = error)
	except KeyboardInterrupt:
		pass
	# Save result
	with open('result.txt', 'w', encoding = 'utf8') as f:
		o.save_txt(f)
	with open('result.json', 'w', encoding = 'utf8') as f:
		o.save_json(f)
