# terraform state dependency walker

Support only when terraform state is stored in aws s3

1. `get-all` command

   - get terraform state files from aws s3 bucket

     - analyze analize terraform state file
       - with go routine

   - save result as a json file
     - (opt) inject result to aws dynamodb

2. `check` command
   - get current backend info
     - user input
   - search dependency
     - from saved json file or aws dynamodb
