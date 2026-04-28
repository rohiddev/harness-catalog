# 🔒 PROTECTED FILE – DO NOT MODIFY

Step 1 — Store workflow YAML in Git (one branch per environment)

┌───────────────┬──────────────────┬──────────────────────────┐
│  Environment  │      Branch      │ Account ID in output URL │                                                                                                                                        
├───────────────┼──────────────────┼──────────────────────────┤                                                                                                                                      
│ P-2 (dev)     │ feature/* or dev │ Dev account ID           │                                                                                                                                        
├───────────────┼──────────────────┼──────────────────────────┤                                                                                                                                      
│ P-1 (staging) │ staging          │ Staging account ID       │                                                                                                                                        
├───────────────┼──────────────────┼──────────────────────────┤                                                                                                                                      
│ P (prod)      │ master           │ Prod account ID          │
└───────────────┴──────────────────┴──────────────────────────┘     

Step 2 — Configure IDP to import from correct branch per environment

IDP > Admin > Git Integrations > Connector — set the branch per environment when registering the workflow.  

===
Step 3 — Proxy endpoints per environment

IDP > Admin > Plugins > Configure Backend Proxies — the fuse-api proxy target URL will differ between P-2 / P-1 / P. Set the correct backend URL per environment here. Your workflow YAML uses       
proxy/fuse-api/... — that path stays the same, only the target resolves differently.

====
Summary of changes needed

┌───────────────────────────────────┬─────────────────────┬──────────────────────────────────────────────────────────────────────┐                                                                   
│               What                │        Where        │                                Action                                │
├───────────────────────────────────┼─────────────────────┼──────────────────────────────────────────────────────────────────────┤
│ Account ID in pipeline            │ Pipeline YAML       │ Replace with <+account.identifier> — one time fix                    │
├───────────────────────────────────┼─────────────────────┼──────────────────────────────────────────────────────────────────────┤                                                                     
│ Account ID in workflow output URL │ Workflow YAML       │ Accept per-branch, OR move to orgIdentifier/projectIdentifier params │
├───────────────────────────────────┼─────────────────────┼──────────────────────────────────────────────────────────────────────┤                                                                     
│ Proxy backend URL                 │ IDP Admin > Proxies │ Set correct target per environment                                   │                                                                   
├───────────────────────────────────┼─────────────────────┼──────────────────────────────────────────────────────────────────────┤                                                                     
│ Workflow import branch            │ IDP Admin > Git     │ Point each env's IDP at its own branch                               │                                                                     
└───────────────────────────────────┴─────────────────────┴──────────────────────────────────────────────────────────────────────┘     

