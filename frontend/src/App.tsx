import { useState, useCallback, useEffect, use } from 'react';
import { ReactFlow, Background, Controls, type Node, type Edge, addEdge, type OnNodesChange, type OnEdgesChange, type OnNodesDelete, applyNodeChanges, applyEdgeChanges, type NodeMouseHandler, type Connection } from '@xyflow/react';
import '@xyflow/react/dist/style.css';
import { QueryClientProvider } from '@tanstack/react-query';
import { queryClient } from './utils/http';

// 引入我們剛剛定義的型別與工具
import type { BaseParams, RetryPolicy, ConditionParams, MoveParams, SleepParams } from './types/schema';
import ConditionNode from './components/nodes/ConditionNode';
import ActionNode from './components/nodes/ActionNode';
import StartNode from './components/nodes/StartNode';
import EndNode from './components/nodes/EndNode';
import WorkflowToolbar from './components/WorkflowToolbar';
import NodeEditorModal from './components/NodeEditorModal';
import { transformBackToReactFlow, transformToDagPayload } from './utils/dagAdapter';
import { useSaveWorkflow } from './hooks/useSaveWorkflow';
import { useFetchWorkflow } from './hooks/useFetchWorkflow';
import { useFetchWorkflowById } from './hooks/useFetchWorkflowById';
import { type WorkflowPayload } from './types/schema';

// --- 註冊 Custom Nodes ---
const nodeTypes = {
  condition: ConditionNode,
  action: ActionNode,
  start: StartNode,
  end: EndNode,
};

function Scheduler() {
  const [nodes, setNodes] = useState<Node[]>([]);
  const [edges, setEdges] = useState<Edge[]>([]);

  // 1. 列表狀態
  const { data: workflows = []} = useFetchWorkflow();
  const [currentWorkflowName, setCurrentWorkflowName] = useState<string>('');
  const [currentWorkflowId, setCurrentWorkflowId] = useState<string>('');

  // 2. 詳細資料狀態 (依賴 currentWorkflowId)
  const { data: workflowDetail, isLoading: isDetailLoading } = useFetchWorkflowById(currentWorkflowId);
  
  // Effect 1: 當列表載入且目前沒有選中時，預設選取第一筆
  useEffect(() => {
    if (workflows.length > 0 && !currentWorkflowId){
      if (workflows[0].workflow_id) {
        setCurrentWorkflowId(workflows[0].workflow_id);
        setCurrentWorkflowName(workflows[0].workflow_name || '');
      }
    }
  }, [workflows, currentWorkflowId])

  // Effect 2: 當詳細資料載入時，更新畫布
  useEffect(() => {
    if (workflowDetail && workflowDetail.nodes) {
        // 使用 Detail 資料更新畫布
        const { nodes: newNodes, edges: newEdges } = transformBackToReactFlow(workflowDetail.nodes);
        setNodes(newNodes);
        setEdges(newEdges);
    } else if (workflowDetail && !workflowDetail.nodes) {
        // 如果該 Workflow 還是空的
        setNodes([]);
        setEdges([]);
    }
  }, [workflowDetail]);

  // Handle Workflow Selection Change
  const handleWorkflowSelect = (id: string) => {
    setCurrentWorkflowId(id);
    // 這裡我們只更新 ID，React Query 的 Hook 會自動幫我們去抓新的資料
    // 並觸發上面的 useEffect 更新畫布
  };

  // --- Modal 狀態 ---
  const [editingNode, setEditingNode] = useState<Node | null>(null);
  const [modalOpen, setModalOpen] = useState(false);

  // React Flow 回呼函式
  const onNodesChange: OnNodesChange = useCallback((changes) => setNodes((nds) => applyNodeChanges(changes, nds)), []);
  const onEdgesChange: OnEdgesChange = useCallback((changes) => setEdges((eds) => applyEdgeChanges(changes, eds)), []);
  const onConnect = useCallback((params: Connection) => setEdges((eds) => addEdge(params, eds)), []);

  const onNodesDelete: OnNodesDelete = useCallback((deleted) => {
    setEdges((eds) => {
      const deletedIds = new Set(deleted.map((node) => node.id));
      return eds.filter((edge) => !deletedIds.has(edge.source) && !deletedIds.has(edge.target));
    });
  }, []);

  // --- 核心：雙擊節點開啟編輯 ---
  const onNodeDoubleClick: NodeMouseHandler = useCallback((_event, node) => {
    if (node.id === 'start' || node.id === 'end') return;
    setEditingNode(node);
    setModalOpen(true);
  }, []);

  // --- 儲存 Modal 的變更 ---
  const handleSaveNodeParams = (newParams: BaseParams | MoveParams | SleepParams | ConditionParams, newRetryPolicy?: RetryPolicy) => {
    if (!editingNode) return;
    
    setNodes((nds) => nds.map((node) => {
      if (node.id === editingNode.id) {
        // 更新節點的 data (UI 會自動更新)
        return { 
          ...node, 
          data: { 
            ...node.data, 
            params: newParams,
            retryPolicy: newRetryPolicy || node.data.retryPolicy
          } 
        };
      }
      return node;
    }));
    setModalOpen(false);
  };

  // --- API Mutation (儲存 Workflow) ---
  const mutation = useSaveWorkflow();
  const handleSaveWorkflow = () => {
    const payload = transformToDagPayload(currentWorkflowId, currentWorkflowName, nodes, edges);
    console.log("Saving Workflow Payload:", JSON.stringify(payload, null, 2));
    // Call mutate inside the handler
    // mutation.mutate(payload, {
    //   onSuccess: () => {
    //      alert('Workflow saved successfully!');
    //   },
    //   onError: (error) => {
    //     alert(`Failed to save: ${error.message}`);
    //   }
    // });

  }

  // --- 新增節點功能 ---
  const handleAddNode = (type: string) => {
    const id = Date.now().toString();
    const isCondition = type === 'Condition';
    
    const newNode: Node = {
      id,
      position: { x: Math.random() * 400, y: Math.random() * 400 },
      data: { 
        label: type,
        activityType: type, 
        params: type === 'Move' ? { x: 0, y: 0 } 
              : type === 'Sleep' ? { duration: 1000 } 
              : type === 'Condition' ? { expression: 'true' }
              : {},
        retryPolicy: {
          maxAttempts: 1,
          initialInterval: 1000,
        }
      },
      type: isCondition ? 'condition' : 'action', // use 'action' instead of 'default'
    };
    setNodes((nds) => [...nds, newNode]);
  };

  // --- 刪除節點功能 ---
  const handleDeleteNode = () => {
    if (!editingNode) return;
    setNodes((nds) => nds.filter((n) => n.id !== editingNode.id));
    setEdges((eds) => eds.filter((e) => e.source !== editingNode.id && e.target !== editingNode.id));
    setModalOpen(false);
    setEditingNode(null);
  };

  return (
    <div className="w-screen h-screen flex flex-col">
      <WorkflowToolbar 
        onAddNode={handleAddNode}
        onSave={handleSaveWorkflow}
        // Pass List and Selection Handler
        workflows={workflows.map((w: WorkflowPayload) => ({ 
            workflow_id: w.workflow_id || '', 
            workflow_name: w.workflow_name || 'Untitled' 
        }))}
        currentWorkflowId={currentWorkflowId}
        onWorkflowSelect={handleWorkflowSelect}
      />

      <div className="flex-1">
        <ReactFlow
          nodes={nodes}
          edges={edges}
          onNodesChange={onNodesChange}
          onEdgesChange={onEdgesChange}
          onConnect={onConnect}
          onNodesDelete={onNodesDelete}
          onNodeDoubleClick={onNodeDoubleClick}
          nodeTypes={nodeTypes}
          fitView
        >
          <Background />
          <Controls />
        </ReactFlow>
      </div>

      <NodeEditorModal
        isOpen={modalOpen}
        node={editingNode}
        onClose={() => setModalOpen(false)}
        onSave={handleSaveNodeParams}
        onDelete={handleDeleteNode}
      />
    </div>
  );
}

// --- 包覆 Provider ---
export default function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <Scheduler />
    </QueryClientProvider>
  );
}
