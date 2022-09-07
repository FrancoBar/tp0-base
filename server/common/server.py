import os
import socket
import logging
import multiprocessing
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
        while self._open:
            client_sock = self.__accept_new_connection()
            if client_sock:
                process = multiprocessing.Process(target=self.__handle_client_connection, args=[client_sock])
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

    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket
        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            pid = os.getpid()
            winners = []

            logging.debug('[{}] Awaiting person record reception'.format(pid))
            personrecords = recv(client_sock)
            logging.debug('[{}] Received {} records'.format(pid, len(personrecords)))

            logging.debug('[{}] Sending back result'.format(pid))
            for p in personrecords:
                if is_winner(p):
                    winners.append(p)
                    send(client_sock, 1)
                else:
                    send(client_sock, 0)

            logging.debug('[{}] Amount of winners: {}'.format(pid, len(winners)))
            persist_winners(winners)
        except InvalidIntentionError:
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