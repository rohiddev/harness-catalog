# 🔒 PROTECTED FILE – DO NOT MODIFY

Harness has a modern UI, YAML-first pipelines, and built-in templates/governance
Harness has built-in Test Intelligence (skips redundant tests based on code change impact), failure strategies, conditional execution, and looping/matrix natively
Harness has OPA-based pipeline policies, approval gates, environment-level RBAC, and audit trails baked in
Smarter CI — Test Intelligence skips unnecessary tests; Jenkins runs everything every time unless you hack it
Faster pipelines — Harness Test Intelligence analyzes what code changed and runs
only the impacted tests, so your pipeline finishes faster even with a large test suite; Jenkins runs the full suite every time regardless of change scope
OPA vs RBAC — Key Difference
RBAC controls who can do what — it's identity-based access (user/role → action → resource)
OPA-based policies control what can be deployed and how — it's rule-based governance on the pipeline itself, regardless of who triggered it
Harness understands your code change graph — only runs tests covering changed classes/methods

Harness lets you define per-stage: retry, rollback, ignore, manual intervention — declaratively in YAML


===
. Deployment Verification
* Harness has built-in CV (Continuous Verification) — monitors APM/logs post-deploy and auto-rolls back on anomaly
* Jenkins deploys and walks away — you find out something broke from your monitoring team
* Beat: Harness closes the feedback loop; Jenkins opens a ticket


==========
Rollback
* Harness rollback is one-click or automatic — goes back to last known good artifact + config
* Jenkins rollback means re-running a previous build manually and hoping state matches
* Beat: MTTR drops significantly in prod incidents

=====
How it Works:
1. You commit desired state (K8s manifests, Helm charts, Terraform) to Git
2. A GitOps agent (Argo CD / Flux) watches that repo
3. Any change in Git → agent automatically syncs it to the cluster
4. If someone manually changes something in the cluster → agent detects drift and corrects it back to what Git says

===

OSS Terraform: you bolt on Sentinel (paid), OPA separately, or write your own CI gates
IaCM: OPA policies are first-class — you enforce cost limits, naming conventions, required tags before apply runs, in the same pipeline
For a banking org, this is non-negotiable — IaCM gives you audit trails without custom tooling

IaCM: managed state backend, versioned, with drift detection built in
Governance, cost estimation, drift detection, and approval workflows are built-in — not bolted on.

====
Dstabase Ops
seamlessly integrates database changes into CI/CD pipelines, automating schema management and governance. It provides full visibility into database changes across environments, ensuring reliable rollbacks and compliance with centralized policies.

Every database change is now governed, not gambled.
It removes the last manual bottleneck in your delivery pipeline

==
What type of change it was. Adding a column? Rollback is straightforward. Dropping a column or table? The rollback command can recreate the structure, but any data written into that column after the migration is gone forever.
Whether rollback scripts exist. Liquibase requires you to explicitly write rollback SQL for certain operations — it doesn't always auto-generate it. If the rollback script wasn't written, the pipeline has nothing to execute.
Whether data was written between deploy and failure. If the app wrote new records in the 30 seconds before the error was detected, rolling back the schema can corrupt or orphan that data.


==
How Harness AI Test Automation Transforms Software Testing:
* 10× Faster Test Creation: The no-code, Generative AI-powered platform enables anyone to generate high-quality test cases within minutes, significantly reducing the time required for test authoring
* 70% Less Test Maintenance: AI-driven test execution and self-healing capabilities minimize the burden of test maintenance, ensuring robust and reliable automated testing with minimal manual intervention.
* 5× Faster Release Cycles: By automating end-to-end testing, Harness AI Test Automation accelerates software delivery, enabling teams to release high-quality applications more frequently and with confidence.

===
Harness Software Engineering Insights (SEI) enables engineering leaders to make data-driven decisions that improve engineering productivity, efficiency, alignment, planning, and execution. It provides actionable insights into software delivery and workflows across teams, processes, and systems to improve software quality, enhance developer experience, and accelerate time to value. Learn how you can use data-led insights to remove bottlenecks and improve productivity. explain this in simple terms.

=====
The 4 DORA Metrics
Metric	What it measures	You want this to be...
Deployment Frequency	How often you ship to production	High (daily or more)
Lead Time for Changes	Time from code commit → production	Low (hours, not weeks)
Change Failure Rate	% of deployments that cause incidents	Low (under 15%)
Mean Time to Restore (MTTR)	How fast you recover from failures	Low (under 1 hour)

DevOps Research and Assessment.

Feature Flags

Release Notes
Harness Feature Flags (FF) is a feature management solution that lets you change your software's functionality without deploying new code. It does this by letting you hide code or behavior without having to ship new versions of the software. A feature flag is like a powerful If statement.



