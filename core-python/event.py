#!/usr/bin/env python3
# -*- coding: utf-8 -*-
# Event Dispatcher

class EventDispatcher():
	'EventDispatcher(): An event dispatcher. Use global variable "event" instead.'

	def __init__(self):
		'Initialize EventDispatcher().'
		self.__events = {}

	def subscribe(self, eventName, callback):
		'Subscribe a callback function to an event.\n\n\teventName: A string of the event name.\n\tcallback: A callback function. Must be defined like "foo(**kwargs)".'
		if eventName not in self.__events:
			self.__events[eventName] = []
		if callback not in self.__events[eventName]:
			self.__events[eventName].append(callback)

	def dispatch(self, eventName, **kwargs):
		'Dispatch an event.\n\n\teventName: The string of the event name.\n\tPass other data with keyword arguments.'
		if eventName not in self.__events:
			return
		for f in self.__events[eventName]:
			f(**kwargs)

event = EventDispatcher()
