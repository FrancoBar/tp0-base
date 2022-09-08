import time
import datetime
import fcntl

""" Winners storage location. """
STORAGE = "./output/winners"


""" Contestant data model. """
class Contestant:
	def __init__(self, first_name, last_name, document, birthdate):
		""" Birthdate must be passed with format: 'YYYY-MM-DD'. """
		self.first_name = first_name
		self.last_name = last_name
		self.document = document
		self.birthdate = datetime.datetime.strptime(birthdate, '%Y-%m-%d')
		
	def __hash__(self):
		return hash((self.first_name, self.last_name, self.document, self.birthdate))

	def __str__(self):
		return '\nFirst Name: ' + self.first_name + '\nLast Name: ' + self.last_name + '\nDocument: ' + str(self.document) + '\nBirthdate: ' +str(self.birthdate)


""" Checks whether a contestant is a winner or not. """
def is_winner(contestant: Contestant) -> bool:
	# Simulate strong computation requirements using a sleep to increase function retention and force concurrency.
	time.sleep(0.001)
	return hash(contestant) % 17 == 0


""" Persist the information of each winner in the STORAGE file."""
def persist_winners(winners: list[Contestant]) -> None:
	with open(STORAGE, 'a+') as file:
		fcntl.flock(file, fcntl.LOCK_EX)
		for winner in winners:
			file.write(f'Full name: {winner.first_name} {winner.last_name} | Document: {winner.document} | Date of Birth: {winner.birthdate.strftime("%d/%m/%Y")}\n')
		fcntl.flock(file, fcntl.LOCK_UN)
