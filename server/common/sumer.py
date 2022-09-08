import time
import signal
import logging

BUSY_SERVER_TOLERANCE = 4

class Sumer:
	"""
	Receives ints from the inputq. If a value is not zero it is added to a total.
	If the value is zero a tuple of total and partial is put into the output queue.
	Partial is a boolean value that indicates if the total is expected to change or not.
	"""
	def __init__(self, queues):
		self._open = True
		self._inputq = queues[0]
		self._outputq = queues[1]

	def sum(self):
		signal.signal(signal.SIGTERM, self.sigterm_handler)
		total = 0
		last_time = time.time()
		partial = False
		try:
			while self._open:
				delta = self._inputq.get()
				if delta == 0:
					actual_time = time.time()
					partial = (actual_time - last_time < BUSY_SERVER_TOLERANCE)
					self._outputq.put((total, partial))
				else:
					last_time = time.time()
					total += delta
		except (ValueError, OSError) as e:
			logging.debug(e)


	def sigterm_handler(self, signum, frame):
		logging.debug('SIGTERM received')
		self._open = False
		self._inputq.close()
		self._outputq.close()
