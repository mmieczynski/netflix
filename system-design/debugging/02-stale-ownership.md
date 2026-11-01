# Stale ownership and overlapping workers

[Debugging index](README.md) · [Preparation index](../README.md)

Practice scenario for the diagram-based distributed debugging round. These
questions are reconstructed preparation material, not a disclosed interview.

Consider two workers that both believe they own the same task:

- Worker A holds a lease.
- It pauses long enough for the lease to expire.
- Worker B acquires ownership.
- Worker A resumes and continues executing.

The follow-up is often effectively: **"You said distributed lock. Why did that
not solve it?"**

Lease expiry changes the coordinator's view. It does not physically stop the
old process.

A deeper answer explains:

- Ownership changes receive a new generation.
- State mutations are conditional on current ownership and state version.
- Stale completion messages cannot regress current state.
- External operations require their own applicable protection.

That last point matters. Rejecting A's eventual database update does not undo
an external action A already performed. A downstream system might enforce
fencing, expected resource versions, or operation idempotency. Those mechanisms
protect different things; name the contract being relied upon.

Useful evidence includes ownership histories, heartbeat gaps, process pauses,
attempt overlap, and downstream calls associated with each generation.
Restarting both workers before collecting these records may make the immediate
symptom disappear while leaving the correctness defect intact.
