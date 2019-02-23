## Create File Sequence Diagram

![Create File Sequence Diagram](Create_File_Sequence_Diagram.png)

### Recreate Diagram

https://www.websequencediagrams.com

```
title Create File Sequence Diagram

User->+CLI: CreateFile(name, S3-URL)
CLI->+S3 Bucket: GetFile()
S3 Bucket->-CLI: File
CLI->+Name Node: CreateFile(name, size)
Name Node->-CLI: DataNode List / Block
CLI->CLI: ChopFileIntoBlocks()
loop per Block
CLI->CLI: DetermineDataNode(DataNodeList)
CLI->+ReceivingDataNode: StoreAndForward(Block, DataNodeList)
ReceivingDataNode->CLI:
end
CLI->-User: Result

loop per DataNodeListEntry Not Storing Block
note over ReceivingDataNode,NextDataNode:
    NextDataNode becomes the ReceivingDataNode
    for next loop iteration
end note
ReceivingDataNode->*-NextDataNode: StoreAndForward(Block, DataNodeList)
end
```
