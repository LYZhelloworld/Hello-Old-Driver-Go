#!/usr/bin/env python3
# -*- coding: utf-8 -*-
# Output result to JSON and TXT files
import json
import config, event

class Output():
	'Output(): Write data to JSON/TXT files or read from JSON file.'

	def __init__(self):
		'Initialize Output().'
		self.resource_list = []
		self.need_check = True

	def add(self, title, url, magnets, error):
		'Add item to the resource list.\n\n\ttitle, url, magnets: Data to be added.\n\terror: Whether it is an error page.\n\nEvent:\n\tAddingData(title, url, magnets)'
		event.event.dispatch('AddingData', title = title, url = url, magnets = magnets)
		self.resource_list.append({'title': title, 'url': url, 'magnets': magnets, 'error': error})
		self.need_check = True
		#self.__remove_duplicated_resources()

	def save_json(self, f):
		'Save resource list to JSON file.\n\nf: A file object opened as "w" mode, with utf-8 encoding.\n\nEvent:\n\tBeforeSavingJSON\n\tAfterSavingJSON'
		self.remove_duplicated_resources()
		event.event.dispatch('BeforeSavingJSON')
		json.dump(self.resource_list, fp = f, indent = config.output['json_indent'], sort_keys = config.output['json_sort_keys'], ensure_ascii = config.output['json_ensure_ascii'])
		event.event.dispatch('AfterSavingJSON')

	def save_txt(self, f):
		'Save resource list to TXT file. Empty pages and error pages will be skipped.\n\nf: A file object opened as "w" mode, with utf-8 encoding.\n\nEvents:\n\tBeforeSavingTXT\n\tAfterSavingTXT'
		self.remove_duplicated_resources()
		event.event.dispatch('BeforeSavingTXT')
		for value in self.resource_list:
			if value['error']: # Skip error pages
				continue
			if len(value['magnets']) < 1: # Skip empty pages
				continue
			print(value['title'], file = f)
			print(value['url'], file = f)
			for magnet in value['magnets']:
				print(config.output['txt_magnet_prefix'] + magnet, file = f)
		event.event.dispatch('AfterSavingTXT')

	def remove_duplicated_resources(self):
		'Remove duplicated resources according to their titles.\n\nEvent:\n\tRemovingDuplicatedResources'
		if self.need_check:
			event.event.dispatch('RemovingDuplicatedResources')
			new_resource_list = []
			for resource in self.resource_list:
				add_flag = True
				for added_resource in new_resource_list:
					if added_resource['title'] == resource['title']:
						add_flag = False
						added_resource['magnets'].extend(resource['magnets'])
						added_resource['magnets'] = list(set(added_resource['magnets']))
						break
				if add_flag:
					new_resource_list.append(resource)
			self.resource_list = new_resource_list
			self.need_check = False

	def load_json(self, f):
		'Load JSON file to resource list.\n\nf: A file object opened as "r" mode, with utf-8 encoding.'
		self.resource_list = json.load(fp = f)
		self.need_check = False
