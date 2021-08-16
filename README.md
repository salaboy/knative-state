# Knative Workflow Controller

This project provides a simple controller to create StateMachine Runners that can execute and track long-running process. 

This controller defines two CRDs:
- StateMachine
- StateMachineRunner

StateMachine CRD defines the StateMachine (definition) structure that we want to track state about and how to transition from a state to the next. 

StateMachineRunner CRD represents aa runtime for a specific StateMachine definition, as it has a reference to the StateMachine definition and can host any number of StateMachineInstances.

## Send Events to a StateMachine

Events are sent to a broker which automatically register triggers for the StateMachine definitions, when a new StateMachineRuner is created.
Find the StateMachineRunner URL by running: 
```bigquery
kubectl get ksvc 
```
or
```bigquery
kubectl get statemachinerunner
```
Create a new StateMachineInstance by sending a request to the StateMachineRunner:

```
curl -X POST http://kservice-buy-tickets-statemachine.default.127.0.0.1.nip.io/statemachines

```

You should get something like this: 
```
{"id":"37718d5c-fe8a-11eb-bad7-ceed770ffae2","current":"","context":null}%
```

You can send CloudEvents to the example StateMachine using `curl` or any CloudEvents SDK


```
curl -X POST -H "Content-Type: application/json" \
  -H "ce-specversion: 1.0" \
  -H "ce-source: curl-command" \
  -H "ce-type: JoinedQueue" \
  -H "ce-id: 123-abc" \
  -H "ce-statemachineid: 5c0eeb45-fe8c-11eb-847b-e65959342a48" \
  -d '{"name":"Salaboy"}' \
  http://broker-ingress.knative-eventing.127.0.0.1.nip.io/default/example-broker

```

```
curl -X POST -H "Content-Type: application/json" \
  -H "ce-specversion: 1.0" \
  -H "ce-source: curl-command" \
  -H "ce-type: ExitedQueue" \
  -H "ce-id: 123-abc" \
  -H "ce-statemachineid: d8bafcd5-eee4-11eb-bdad-820448832986" \
  -d '{"name":"Salaboy"}' \
  http://broker-ingress.knative-eventing.127.0.0.1.nip.io/default/example-broker
```

## CRDS

StateMachine (definition)
```
apiVersion: flow.knative.dev/v1
kind: StateMachine
metadata:
  name: buy-tickets-statemachine
spec:
  stateMachine:
    version: "1.0.0"
    states:
      "":
        events:
          JoinedQueue: InQueue
      InQueue:
        events:
          ExitedQueue: BuyingTickets
          AbandonedQueue: NoTicketsPurchased
      BuyingTickets:
        events:
          TicketsReserved: PayingTickets
      PayingTickets:
        events:
          PaymentRequested: WaitingForPaymentAuth
      WaitingForPaymentAuth:
        events:
          PaymentAuthorized: TicketsPurchased
      TicketsPurchased:
        events: { }
      NoTicketsPurchased:
        events: { }

```

StateMachineRunner

```
apiVersion: flow.knative.dev/v1
kind: StateMachineRunner
metadata:
  name: buy-tickets-runner
spec:
  stateMachineRef: buy-tickets-buy-tickets-statemachine
  sink: http://sockeye.default.svc.cluster.local/
```
