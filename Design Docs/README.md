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
Name Node->-CLI: DataNode List / BlockID
CLI->CLI: ChopFileIntoBlocks()

loop per BlockID
CLI->CLI: DetermineDataNode(DataNodeList)
CLI->+ReceivingDataNode: StoreAndForward(Block, BlockID, DataNodeList)
ReceivingDataNode->CLI:
end

CLI->-User: Result

loop per DataNodeListEntry Not Storing BlockID
note over ReceivingDataNode,NextDataNode:
    NextDataNode becomes the ReceivingDataNode
    for next loop iteration
end note
ReceivingDataNode->*-NextDataNode: StoreAndForward(Block, BlockID, DataNodeList)
end
```

## Get File Sequence Diagram

![Get File Sequence Diagram](Get_File_Sequence_Diagram.png)

```
title Get File Sequence Diagram

User->+CLI: GetFile(name, saveLocation)
CLI->+Name Node: GetFile(name, size)
Name Node->-CLI: DataNode List / BlockID

loop per BlockID
CLI->CLI: DetermineDataNode(DataNodeList)
CLI->+DeterminedDataNode: GetBlock(BlockID)
DeterminedDataNode->CLI: Block
end

CLI->CLI: ConcatenateBlocks(blocks)
CLI->CLI: SaveFile(data, saveLocation)
CLI->-User: Result
```
