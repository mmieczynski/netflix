# Lost acknowledgment after a successful operation

[Debugging index](README.md) · [Preparation index](../README.md)

Practice scenario for the diagram-based distributed debugging round. These
questions are reconstructed preparation material, not a disclosed interview.

The sequence is:

1. A worker requests a deployment.
2. The deployment platform accepts it.
3. The response is lost.
4. The worker times out.
5. The work item becomes eligible for retry.

The key reasoning is that the worker has an **unknown outcome**. A timeout does
not establish that the deployment failed.

| Interviewer asks | The answer should establish |
|---|---|
| **"Would you retry?"** | Whether retrying the same logical operation is safe under the downstream contract. |
| **"We store an idempotency key in our database. Is that enough?"** | Not if the downstream effect can commit independently and the record does not atomically protect it. |
| **"What happens if we record success before calling downstream?"** | A crash can leave a false success without the effect. |
| **"What if we record success afterward?"** | A crash can leave a completed effect with no local success record. |
| **"The external API has no idempotency support."** | Look for a stable external identity, conditional mutation, or authoritative lookup. Otherwise, some ambiguity may require conservative reconciliation or intervention. |
| **"How would you investigate an incident?"** | Correlate logical operation, attempts, external references, request times, and authoritative downstream state. |
| **"How do you recover safely?"** | Reconcile uncertain operations first and retry only under a defined duplicate-safety contract. |

AWS's discussion of idempotent APIs is a strong source for the distinction
between repeated delivery and repeated intent.
[Making retries safe with idempotent APIs](https://builder.aws.com/content/3Ev0BENTyBr0XxzRk5FDZzgNYos/making-retries-safe-with-idempotent-apis)
