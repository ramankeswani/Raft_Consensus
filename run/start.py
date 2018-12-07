#!/usr/bin/python

# All SSH libraries for Python are junk (2011-10-13).
# Too low-level (libssh2), too buggy (paramiko), too complicated
# (both), too poor in features (no use of the agent, for instance)

# Here is the right solution today:

import subprocess
import sys
import os
import threading

com = [line.rstrip() for line in open('run_arg48.txt')]

def run_raft(cmd):
    print(cmd)
    os.system(cmd)


if __name__ == "__main__":
    t = []
    with open("ip48.txt") as f:
        for i, line in enumerate(f):
            #do something with data
            HOST=line.rstrip()
            # Ports are handled in ~/.ssh/config since we use OpenSSH
            #COMMAND=sys.argv[1]
            print(HOST, com[i])
            #print("gcloud compute ssh {0} -- '{1} {2} &; ls'".format(HOST, COMMAND, com[i]) )
            #os.system( "gcloud compute ssh {0} -- '{1} {2} > /dev/null 2>&1 &'".format(HOST, COMMAND, com[i] ) );
            
            cmd = "gcloud compute ssh {0} -- 'nohup ./test/run_process.sh {1} {2} {3}' ".format(HOST, com[i], com[i + 16], com[i+32])
            #run_raft(cmd)
            t_ = threading.Thread(target=run_raft, args=(cmd,))
            t_.start()
            t.append(t_)

            #os.system( cmd ) #.format(HOST, COMMAND, com[i] ) );
            #os.system( "scp GCPSetup.sh ramankeswani{0}:~/".format(line) );
            #ssh = subprocess.Popen(["gcloud compute ssh", "%s" % HOST, COMMAND],
            #                   shell=False,
            #                   stdout=subprocess.PIPE,
            #                   stderr=subprocess.PIPE)
            #result = ssh.stdout.readlines()
            #if result == []:
            #    error = ssh.stderr.readlines()
            #    print >>sys.stderr, "ERROR: %s" % error
            #else:
            #    print result
    for t_ in t:
        t_.join()
