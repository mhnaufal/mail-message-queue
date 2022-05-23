# Not a Project, Unfortunatelly

Here, I try to:

1. ✅ Run message queue using Go
2. ✅ Send email using Go
3. ❌ Send email asynchronously using message queue

---

> producer == sender

> consumer == receiver == worker

> queue == broker

> message == job == task

---

### Receiver + Sender

Sender send a message and put it inside queue, then the receiver will receive the message by taking out the message from the queue.
If there are multiple receivers, the message will be distributed equally (Round-Robin dispatching)

### Worker Queue

Prevent from resource-intensive task by having a queue to schedule the task for later.
**Task** encapsulated into _message_, and send into queue. A _worker_ process will pop up the task and execute the **job**.

If there is a job that times to be done, and suddenly the worker dies and cannot while doing the job, the job will be mark as done
at the queue, even though it is actually not done yet. To prevent this, use the `ack: true` on the producer and set the `message.Ack(false)` on the worker so that when worker die and the job is not finish yet, it will not be mark as done by the queue

RabbitMQ doesn't know the level of complexity of the jobs, so it will always distribute the jobs evenly across the workers (fair dispatch).
Use the `ch.Qos` to define the **prefatch count** a.k.a how many jobs the worker can handle at the same time. If set to 1, it means only
1 job can the worker do at that time. Then the job will be given to the others workers available. ⛔ Be carefull if there are no free workers, it can make the queue error.

When the RabbitMQ server stops, it will forget all of the queues and messages. To prevent this, we need to set the durability of the queue
and the message.
