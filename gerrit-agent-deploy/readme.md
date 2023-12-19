kubectl create secret generic gerrit-agent-secret \
--from-literal=GERRIT_USER=gerrit \
--from-literal=GERRIT_PWD=gerrit \
--from-literal=ES_USER=elastic \
--from-literal=ES_PWD=elastic \
-n ems-test