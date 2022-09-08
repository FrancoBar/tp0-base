import os
import signal
import socket
import logging
from asyncio import IncompleteReadError
from .utils import *
from .transmition import *
from .serialize import *

INTENTION_ASK_WINNER = 0
INTENTION_ASK_AMOUNT = 1

class Server:
    def __init__(self, port, listen_backlog):
        self._open = True
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)

    def run(self):
        """
        Dummy Server loop
        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """

        while self._open:
            client_sock = self.__accept_new_connection()
            if client_sock:
                self.__handle_client_connection(client_sock)
        logging.info('Shutting down...')

    def sigterm_handler(self, signum, frame):
        logging.debug('SIGTERM received')
        self._open = False
        logging.debug('Closing socket')
        self._server_socket.close()

    def __ask_winner(self, client_sock):
        pid = os.getpid()
        winners = []
        logging.debug('[{}] Awaiting person record reception'.format(pid))
        personrecords = recv_vector(client_sock, recv_person_record)
        logging.debug('[{}] Received {} records'.format(pid, len(personrecords)))

        logging.debug('[{}] Sending back result'.format(pid))
        for p in personrecords:
            has_won = is_winner(p)
            if has_won:
                winners.append(p)
            send_bool(client_sock, has_won)

        winners_amount = len(winners)
        logging.debug('[{}] Amount of winners: {}'.format(pid, winners_amount))

    def __handle_client_connection(self, client_sock):
            """
            Read message from a specific client socket and closes the socket
            If a problem arises in the communication with the client, the
            client socket will also be closed
            """
            try:
                pid = os.getpid()
                while self._open:
                    intention = recv_uint32(client_sock)
                    if intention == INTENTION_ASK_WINNER:
                        self.__ask_winner(client_sock)
                    else:
                        logging.info('[{}] Error: Client sent an invalid intention'.format(pid))

            except IncompleteReadError as e:
                logging.debug("[{}] {}".format(pid, e))
            except (OSError, ValueError) as e:
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