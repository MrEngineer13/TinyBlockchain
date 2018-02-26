# TinyBlockchain
This project was a port of [snakecoin](https://medium.com/crypto-currently/lets-make-the-tiniest-blockchain-bigger-ac360a328f4d) to Go.

The idea is to have a super simple server because the block chain 

# Usage
Start everything:

    go run web.go

Add transactions:

    curl "localhost:9090/txion" \
         -H "Content-Type: application/json" \
         -d '{"from": "akjflw", "to":"fjlakdj", "amount": 3}'
      
Navigate in a browerser 

    http://localhost:9090/ -> Generates blocks to add
    http://localhost:9090/mine -> Mines a new block and adds it to the chain
    http://localhost:9090/blocks -> Makes sure you have the longest blockchain then prints it
