import { useState, useCallback, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useLocation } from 'react-router-dom';
import { ReactFlow, Background, Controls, type Node, type Edge, addEdge, type OnNodesChange, type OnEdgesChange, type OnNodesDelete, applyNodeChanges, applyEdgeChanges, type NodeMouseHandler, type Connection } from '@xyflow/react';
import '@xyflow/react/dist/style.css';

import type { BaseParams, MoveParams, SleepParams } from '@/types/workflows';
import ActionNode from '@/components/nodes/ActionNode';
import StartNode from '@/components/nodes/StartNode';
import EndNode from '@/components/nodes/EndNode';
import NodeEditorModal from '@/components/NodeEditorModal';
import { WorkflowStatusBadge } from '@/components/shared/WorkflowStatusBadge';
import { transformBackToReactFlow, transformToDagPayload } from '@/utils/dagAdapter';
import { useSaveWorkflow } from '@/hooks/useSaveWorkflow';
import { useFetchWorkflowById } from '@/hooks/useFetchWorkflows';
import { useWorkflowMonitor } from '@/hooks/useWorkflowMonitor';
import { useFetchActivitiesDef } from '@/hooks/useFetchActivitiesDef';
import { useTriggerWorkflow, usePauseWorkflow, useResumeWorkflow } from '@/hooks/useControlWorkflow';
import { Button } from '@/components/ui/button';
import { ArrowLeft, Save, Plus, Pause, Play, Pencil } from 'lucide-react';
import type { ActivityDefinition } from '@/types/activities';

// --- Register Custom Nodes ---
const nodeTypes = {
  action: ActionNode,
  start: StartNode,
  end: EndNode,
};

const defaultNodes = [{
  id: 'start',
  type: 'start',
  position: { x: 0, y: 0 },
  data: { label: 'Start', activityType: "Start"},
},{
  id: 'end',
  type: 'end',
  position: { x: 1000, y: 0 },
  data: { label: 'End', activityType: "End"},
}];

export default function Editor() {
  const { workflowId } = useParams<{ workflowId: string }>();

  // state passed by useLocation in react-router
  const location = useLocation();
  const state = location.state as { operation: string; workflowName?: string };
  const isCreateMode = state.operation === 'create';

  // 編輯模式 vs 執行模式
  const [isEditing, setIsEditing] = useState(true);
  
  // Fetch workflow detail from backend
  const { data: workflowDetail, isLoading } = useFetchWorkflowById(workflowId || '', { enabled: !isCreateMode});
  const workflowName = workflowDetail?.workflow_name || state.workflowName;

  // Fetch activity definitions
  const { data: activitiesDef } = useFetchActivitiesDef();

  // React Router navigation
  const navigate = useNavigate();
  
  // --- React Flow State ---
  const [nodes, setNodes] = useState<Node[]>(defaultNodes);
  const [edges, setEdges] = useState<Edge[]>([]);

  // Effect: When workflow detail is loaded, update the canvas
  useEffect(() => {
    if (!workflowDetail || !workflowDetail.nodes) return;

    const { nodes: newNodes, edges: newEdges } = transformBackToReactFlow(
      workflowDetail.nodes,
      activitiesDef
    );
    setNodes(newNodes);
    setEdges(newEdges);
  }, [workflowDetail, activitiesDef]);

  // --- Workflow 執行狀態監控 ---
  const { setIsMonitoring, workflowStatus, currentNode, currentStep } = useWorkflowMonitor(workflowId!);

  // 當 workflow 開始執行時，自動切換到執行模式
  useEffect(() => {
    if (workflowStatus === 'Running' || workflowStatus === 'Paused') {
      setIsEditing(false);
    }
  }, [workflowStatus]);
  
  // Effect: 當 currentNode 改變時，更新節點的 running 狀態
  useEffect(() => {
    setNodes((nds) =>
      nds.map((node) => ({
        ...node,
        data: {
          ...node.data,
          isRunning: node.id === currentNode,
        },
      }))
    );
  }, [currentNode]);

  // --- React Flow Node Modal state ---
  const [editingNode, setEditingNode] = useState<Node | null>(null);
  const [modalOpen, setModalOpen] = useState(false);

  // React Flow callback function
  const onNodesChange: OnNodesChange = useCallback((changes) => setNodes((nds) => applyNodeChanges(changes, nds)), []);
  const onEdgesChange: OnEdgesChange = useCallback((changes) => setEdges((eds) => applyEdgeChanges(changes, eds)), []);
  const onConnect = useCallback((params: Connection) => setEdges((eds) => addEdge(params, eds)), []);
  const onNodesDelete: OnNodesDelete = useCallback((deleted) => {
    setEdges((eds) => {
      const deletedIds = new Set(deleted.map((node) => node.id));
      return eds.filter((edge) => !deletedIds.has(edge.source) && !deletedIds.has(edge.target));
    });
  }, []);

  // --- Core: Double-click node to open editor ---
  const onNodeDoubleClick: NodeMouseHandler = useCallback((_event, node) => {
    if (!isEditing) return; // 執行模式下不能編輯
    if (node.id === 'start' || node.id === 'end') return;
    setEditingNode(node);
    setModalOpen(true);
  }, [isEditing]);

  // --- Save changes from Modal ---
  const handleSaveNodeParams = (newParams: BaseParams | MoveParams | SleepParams) => {
    if (!editingNode) return;
    
    setNodes((nds) => nds.map((node) => {
      if (node.id === editingNode.id) {
        return { 
          ...node, 
          data: { 
            ...node.data, 
            params: newParams,
          } 
        };
      }
      return node;
    }));
    setModalOpen(false);
  };

  // --- API Mutation (儲存 Workflow) ---
  const saveMutation = useSaveWorkflow();
  const handleSaveWorkflow = () => {
    if (!workflowId) return;
    const payload = transformToDagPayload(nodes, edges);
    saveMutation.mutate({
      workflow_id: workflowId,
      workflow_name: workflowName,
      nodes: payload,
    }, {
      onSuccess: () => {
        alert("Workflow saved successfully.");
      },
      onError: (error) => {
        alert(`Failed to save: ${error.message}`);
      }
    });
  };


  const triggerMutation = useTriggerWorkflow();
  const handleTriggerWorkflow = () => {
    if (!workflowId) return;
    triggerMutation.mutate(workflowId, {
      onSuccess: () => {
        setIsMonitoring(true);
      },
      onError: (error) => {
        console.log(`Failed to trigger: ${error.message}`);
        setIsMonitoring(false);
      }
    });
  };

  const pauseWorkflow = usePauseWorkflow();
  const handlePauseWorkflow = () => {
    if (!workflowId) return;
    pauseWorkflow.mutate(workflowId, {
      onSuccess: () => {
        alert('Workflow paused successfully!');
      },
      onError: (error) => {
        alert(`Failed to pause: ${error.message}`);
      }
    });
  };

  const resumeWorkflow = useResumeWorkflow();
  const handleResumeWorkflow = () => {
    if (!workflowId) return;
    resumeWorkflow.mutate(workflowId, {
      onSuccess: () => {
        alert('Workflow resumed successfully!');
      },
      onError: (error) => {
        alert(`Failed to resume: ${error.message}`);
      }
    });
  };

  // --- 新增節點功能 ---
  const handleAddNode = useCallback((activity: ActivityDefinition) => {
    const id = Date.now().toString();

    const defaultParams: Record<string, any> = {};
    if (activity.input_schema?.properties) {
      Object.entries(activity.input_schema.properties).forEach(([key, prop]) => {
        defaultParams[key] = prop.default ?? (prop.type === 'number' ? 0 : prop.type === 'string' ? '' : null);
      });
    }

    const newNode: Node = {
      id,
      position: { x: Math.random() * 400, y: Math.random() * 400 },
      data: { 
        label: activity.name,
        activityType: activity.activity_type,
        inputSchema: activity.input_schema,
        params: defaultParams,
      },
      type: activity.node_type,
    };
    setNodes((nds) => [...nds, newNode]);
  }, []);

  // --- 刪除節點功能 ---
  const handleDeleteNode = () => {
    if (!editingNode) return;
    setNodes((nds) => nds.filter((n) => n.id !== editingNode.id));
    setEdges((eds) => eds.filter((e) => e.source !== editingNode.id && e.target !== editingNode.id));
    setModalOpen(false);
    setEditingNode(null);
  };

  if (isLoading) {
    return (
      <div className="w-screen h-screen flex items-center justify-center">
        <div className="text-muted-foreground">Loading workflow...</div>
      </div>
    );
  }

  const isRunning = workflowStatus === 'Running';
  const isPaused = workflowStatus === 'Paused';

  return (
    <div className="w-screen h-screen flex flex-col">
      {/* Toolbar */}
      {/* Toolbar - 只在大螢幕顯示 */}
      <header className="border-b bg-card px-4 py-3 hidden lg:block">
        {/* Desktop Layout */}
        <div className="flex items-center justify-between">
          {/* Left: Back + Title + Mode Toggle */}
          <div className="flex items-center gap-4">
            <Button variant="ghost" size="sm" onClick={() => navigate('/')}>
              <ArrowLeft className="h-4 w-4 mr-2" />
              Back
            </Button>
            <div className="h-6 w-px bg-border" />
            <div>
              <h1 className="font-semibold">{workflowName}</h1>
              <WorkflowStatusBadge status={workflowStatus} currentStep={currentStep || ""} />
            </div>
            {/* 模式切換按鈕 */}
            <div className="h-6 w-px bg-border" />
            <div className="flex items-center gap-1 bg-muted rounded-lg p-1">
              <Button
                variant={isEditing ? 'default' : 'ghost'}
                size="sm"
                onClick={() => setIsEditing(true)}
                disabled={isRunning || isPaused}
              >
                <Pencil className="h-3 w-3 mr-1" />
                Edit
              </Button>
              <Button
                variant={!isEditing ? 'default' : 'ghost'}
                size="sm"
                onClick={() => setIsEditing(false)}
              >
                <Play className="h-3 w-3 mr-1" />
                Run
              </Button>
            </div>
          </div>

          {/* Center: Add Node Buttons - 只在編輯模式顯示 */}
          {isEditing ? (
            <div className="flex items-center gap-2">
              {activitiesDef?.map((activity) => (
                <Button
                  key={activity.activity_type}
                  variant="outline"
                  size="sm"
                  onClick={() => handleAddNode(activity)}
                >
                  <Plus className="h-3 w-3 mr-1" />
                  {activity.name}
                </Button>
              ))}
            </div>
          ) : (
            <div className="text-sm text-muted-foreground">
              Switch to Edit mode to modify workflow
            </div>
          )}

          {/* Right: Actions */}
          <div className="flex items-center gap-2">
            {/* Save - 只在編輯模式顯示 */}
            {isEditing && (
              <>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={handleSaveWorkflow}
                >
                  <Save className="h-4 w-4 mr-2" />
                  Save
                </Button>
              </>
            )}

            {/* Run/Pause/Resume - 只在執行模式顯示 */}
            {!isEditing && (
              <>
                {isRunning || isPaused ? (
                  <Button
                    variant={isPaused ? 'default' : 'secondary'}
                    size="sm"
                    onClick={isPaused ? handleResumeWorkflow : handlePauseWorkflow}
                  >
                    {isPaused ? (
                      <>
                        <Play className="h-4 w-4 mr-2" />
                        Resume
                      </>
                    ) : (
                      <>
                        <Pause className="h-4 w-4 mr-2" />
                        Pause
                      </>
                    )}
                  </Button>
                ) : (
                  <Button
                    variant="default"
                    size="sm"
                    onClick={handleTriggerWorkflow}
                    disabled={triggerMutation.isPending}
                  >
                    <Play className="h-4 w-4 mr-2" />
                    {triggerMutation.isPending ? 'Starting...' : 'Run'}
                  </Button>
                )}
              </>
            )}
          </div>
        </div>
      </header>

      {/* Mobile: 只顯示返回按鈕，toolbar 隱藏 */}
      <header className="border-b bg-card px-3 py-2 flex items-center justify-between lg:hidden">
        <Button variant="ghost" size="sm" onClick={() => navigate('/')}>
          <ArrowLeft className="h-4 w-4" />
        </Button>
        <div className="text-center flex-1">
          <h1 className="font-semibold text-sm truncate">{workflowName}</h1>
          <WorkflowStatusBadge status={workflowStatus} currentStep={currentStep || ""} />
        </div>
        <div className="w-8" /> {/* Spacer for balance */}
      </header>

      {/* React Flow Canvas */}
      <div className="flex-1 bg-slate-100">
        <ReactFlow
          nodes={nodes}
          edges={edges}
          onNodesChange={isEditing ? onNodesChange : undefined}
          onEdgesChange={isEditing ? onEdgesChange : undefined}
          onConnect={isEditing ? onConnect : undefined}
          onNodesDelete={isEditing ? onNodesDelete : undefined}
          onNodeDoubleClick={onNodeDoubleClick}
          nodeTypes={nodeTypes}
          nodesDraggable={isEditing}
          nodesConnectable={isEditing}
          elementsSelectable={isEditing}
          fitView
        >
          <Background />
          <Controls />
        </ReactFlow>
      </div>

      {/* Node Editor Modal */}
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
