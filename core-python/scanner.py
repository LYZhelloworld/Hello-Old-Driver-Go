#!/usr/bin/env python3
# -*- coding: utf-8 -*-
# Crawling pages and return result
import requests, re
import config, event

class Scanner():
	'Scanner(domain): Get pages according to article ID.\n\nUsage:\n\tfor url, result_text in scanner_instance:\n\t\t...'

	def __init__(self, domain):
		'Initialize Scanner(domain).\n\ndomain: The domain used in URL like "http://domain/".'
		self.__viewed_urls = []
		self.__session = requests.Session()
		self.__session.headers.update({'Cookie': config.scanner['cookie']})
		self.__resource_list = []
		self.domain = domain
		self.maximum = config.scanner['maximum_id']

	def __iter__(self):
		'Implement iter(self).\n\nEvent:\n\tStartScanning'
		event.event.dispatch('StartScanning')
		#return self.__scan_page(self, self.root_url)
		self.__current_id = 1
		return self

	def __next__(self):
		'Implement next(self).\n\nEvent:\n\tScanningPage(url)'
		if self.__current_id > self.maximum: # Maximum ID reached
			raise StopIteration

		url = '{0}://{1}/wp/{2}.html'.format(config.scanner['protocol'], self.domain, self.__current_id)
		while url in self.__viewed_urls: # Skip visited urls
			self.__current_id += 1
			url = '{0}://{1}/wp/{2}.html'.format(config.scanner['protocol'], self.domain, self.__current_id)

		# Retrieve page
		result = self.__session.get(url, timeout = 60)
		if result.status_code != 200: # Error
			self.__current_id += 1
			self.__viewed_urls.append(url)
			return (url, None)

		event.event.dispatch('ScanningPage', url = url)
		result_text = result.content.decode()

		# Next
		self.__current_id += 1
		self.__viewed_urls.append(url)
		return (url, result_text)

	def next(self): # Python 2 compatible
		'Implement iter(self).'
		return self.__next__

	def load_visited_urls(self, resource_list):
		'Load visited URLs from a resource list.'
		for value in resource_list:
			self.__viewed_urls.append(value['url'])
