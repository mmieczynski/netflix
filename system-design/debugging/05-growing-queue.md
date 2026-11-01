# Queue growth that does not improve with more workers

[Debugging index](README.md) · [Preparation index](../README.md)

Practice scenario for the diagram-based distributed debugging round. These
questions are reconstructed preparation material, not a disclosed interview.

The interviewer may give a queue between several components and ask what could
go wrong.

Distinguish:

- Producers generating more work.
- Consumers processing each item more slowly.
- Retries creating additional work.
- One partition or tenant receiving disproportionate traffic.
- A poison message repeatedly failing.
- Consumers waiting on a saturated downstream dependency.
- Consumers processing successfully but failing to commit progress.

Useful measurements include arrival rate, completion rate, oldest-item age,
per-partition lag, attempt counts, processing-time distributions, and downstream
saturation.

For an illustrative calculation:

> At 1,200 arrivals per second and 1,000 completions per second, backlog grows
> by 200 items per second.

Once capacity exceeds arrivals, the backlog drains at the **difference**
between those rates. Adding workers does not help if they all contend for the
same database bottleneck; it may make that bottleneck worse.

A good mitigation explains what happens to delayed work, whether ordering is
preserved, and whether anything is dropped. Moving a message to a dead-letter
queue is only useful if ownership, diagnosis, and replay are defined.
