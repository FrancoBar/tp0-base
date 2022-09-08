import os
import socket
import logging
import multiprocessing
import time
from .utils import *
from .transmition import *
from asyncio import IncompleteReadError

class Server:
    def __init__(self, port, listen_backlog):
        self._open = True
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)

    def run(self):
        """
        Server that accept a new connections and establishes a
        communication with a client. Each connection is handled by
        a separate process.
        """
        winners_queue = (multiprocessing.Queue(), multiprocessing.Queue()) 
        process = multiprocessing.Process(target=self.__sum_winners, args=[winners_queue])
        process.start()
        
        while self._open:
            client_sock = self.__accept_new_connection()
            if client_sock:
                process = multiprocessing.Process(target=self.__handle_client_connection, args=[client_sock, winners_queue])
                process.daemon = True
                process.start()
        
        for process in multiprocessing.active_children():
            logging.debug("Terminating process %r", process)
            process.terminate()
            process.join()
        logging.info('Shutting down...')

    def sigterm_handler(self, signum, frame):
        logging.debug('SIGTERM received')
        self._open = False
        logging.debug('Closing socket')
        self._server_socket.close()

    def __sum_winners(self, winners_queue):
        inputq, outputq = winners_queue
        total = 0
        last_time = time.time()
        partial = False
        while self._open:
            delta = inputq.get()
            if delta == 0:
                actual_time = time.time()
                partial = (actual_time - last_time < 4)
                outputq.put((total, partial))
            else:
                last_time = time.time()
                total += delta
            

    def __ask_winner(self, client_sock, winners_queue):
        pid = os.getpid()
        winners = []
        logging.debug('[{}] Awaiting person record reception'.format(pid))
        personrecords = recv_vector(client_sock, recv_person_record)
        logging.debug('[{}] Received {} records'.format(pid, len(personrecords)))

        logging.debug('[{}] Sending back result'.format(pid))
        for p in personrecords:
            if is_winner(p):
                winners.append(p)
                send(client_sock, 1)
            else:
                send(client_sock, 0)

        logging.debug('[{}] Amount of winners: {}'.format(pid, len(winners)))
        winners_queue[0].put(len(winners))
        persist_winners(winners)

    def __ask_amount(self, client_sock, winners_queue):
        winners_queue[0].put(0)
        total, partial = winners_queue[1].get()
        send(client_sock, total)
        send(client_sock, 1 if partial else 0)

    def __handle_client_connection(self, client_sock, winners_queue):
        """
        Read message from a specific client socket and closes the socket
        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            pid = os.getpid()


            while self._open:
                intention = recv_intention(client_sock)
                if intention == INTENTION_ASK_WINNER:
                    self.__ask_winner(client_sock, winners_queue)
                elif intention == INTENTION_ASK_AMOUNT:
                    self.__ask_amount(client_sock, winners_queue)
                else:
                    logging.info('[{}] Error: Client sent an invalid intention'.format(pid))

        except (OSError, IncompleteReadError) as e:
            logging.info("[{}] {}".format(pid, e))
        finally:
            client_sock.close()

    def __accept_new_connection(self):
        """
        Accept new connections
        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """

        # Connection arrived
        logging.info("Proceed to accept new connections")
        try:
            c, addr = self._server_socket.accept()
            logging.info('Got connection from {}'.format(addr))
        except OSError:
            if self._open:
                logging.info("Error while reading socket {}".format(self._server_socket))
            c = None
        return c