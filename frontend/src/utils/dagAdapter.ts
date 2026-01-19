import type { Edge, Node } from '@xyflow/react';
import type { NodeInfo } from '../types/workflows';

export const transformToDagPayload = (
  nodes: Node[],
  edges: Edge[]
): Record<string, NodeInfo> => {
  const payloadNodes: Record<string, NodeInfo> = {};

  // 1. 初始化所有節點
  nodes.forEach((node) => {
    payloadNodes[node.id] = {
      id: node.id,
      type: node.data.activityType,
      params: node.data.params,
      transitions: {},
    };
  });

  // 2. 根據 Edges 建立關聯
  edges.forEach((edge) => {
    const sourceNode = payloadNodes[edge.source];
    const targetId = edge.target;

    if (!sourceNode) return;

    // 根據 Handle ID 對應到正確的 transition slot
    // Action Node handles: 'success' (or null), 'failure'
    // Condition Node handles: 'true', 'false'
    const handleId = edge.sourceHandle;

    if (handleId === 'failure') {
      sourceNode.transitions.failure = targetId;
    } else {
      // default / success handle
      sourceNode.transitions.next = targetId;
    }
  });
  
  return payloadNodes;
};


export const transformBackToReactFlow = (
  backendNodes: Record<string, NodeInfo>
): { nodes: Node[]; edges: Edge[] } => {

  const nodes: Node[] = [];
  const edges: Edge[] = [];

  // 簡單的自動佈局參數 (因為後端目前沒有存 x, y)
  let gridX = 250;
  let gridY = 50;
  const colWidth = 250;
  const rowHeight = 150;
  let colIndex = 0;
  const maxCols = 3;

  Object.values(backendNodes).forEach((node) => {
    
    // 1. 決定 Node Type 與 Position
    let type = "action";
    const activityType = node.type.toLowerCase();
    let position = {x: 0, y: 0};

    if (activityType === "start"){
      type = "start";
      position = { x: 50, y: 300 };
    } else if (activityType === "end"){
      type = "end";
      position = { x: 1000, y: 300 };
    } else {
      position = { x: gridX, y: gridY };

      // 更新 Grid 指標
      gridX += colWidth;
      colIndex++;
      if (colIndex >= maxCols) {
        colIndex = 0;
        gridX = 250; // 重置 X (保持在 Start 右側)
        gridY += rowHeight;
      }
    
    }

    // 2. 計算位置 (Grid Layout)
    const newNode: Node = {
      id: node.id,
      type: type,
      position: position,
      data: {
        label: node.type,
        activityType: node.type,
        params: node.params || {},
      }
    };
    nodes.push(newNode);

    // 2. 建立 Edges (還原連線關係)
    // 必須檢查每個 transition 是否存在，並建立對應的 Edge
    const { transitions } = node;

    const addEdge = (targetId: string | undefined, handleId?: string, label?: string, stroke?: string) =>{
        if (!targetId) return;
        edges.push({
            id: `e-${node.id}-${targetId}-${handleId || 'default'}`,
            source: node.id,
            target: targetId,
            sourceHandle: handleId,
            label: label,
            animated: false,
            style: stroke ? { stroke } : undefined,
        });
    }   

    addEdge(transitions.next);
    addEdge(transitions.failure, 'failure', undefined, '#ff4d4f'); // 紅線代表失敗
  })

  return {nodes, edges};
}