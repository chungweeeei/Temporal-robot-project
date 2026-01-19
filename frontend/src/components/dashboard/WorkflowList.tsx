import { useMemo, useState, useEffect } from "react";
import { WorkflowCard } from "./WorkflowCard";
import { useFetchWorkflow } from "@/hooks/useFetchWorkflow";
import { useFetchWorkflowStatus } from "@/hooks/useFetchWorkflowStatus";
import { useTriggerWorkflow } from "@/hooks/useTriggerWorkflow";
import { usePauseWorkflow, useResumeWorkflow } from "@/hooks/useControlWorkflow";
import { useDeleteWorkflow } from "@/hooks/useDeleteWorkflow";
import { useWorkflowStore } from "@/store/useWorkflowStore";
import type { WorkflowPayload, WorkflowStatus } from "@/types/schema";
import { Workflow } from "lucide-react";

export function WorkflowList() {
  const { data: workflows = [], isLoading } = useFetchWorkflow();
  const { activeWorkflowId, setActiveWorkflowId } = useWorkflowStore();
  
  // 只對 running 的 workflow 進行輪詢
  const { data: statusData } = useFetchWorkflowStatus(
    activeWorkflowId || "", 
    !!activeWorkflowId
  );

  const triggerWorkflow = useTriggerWorkflow();
  const pauseWorkflow = usePauseWorkflow();
  const resumeWorkflow = useResumeWorkflow();
  const deleteWorkflow = useDeleteWorkflow();

  // 追蹤每個 workflow 的狀態
  const [workflowStatuses, setWorkflowStatuses] = useState<Record<string, { status: WorkflowStatus; currentStep?: string }>>({});

  // 當 statusData 更新時，更新對應 workflow 的狀態
  useEffect(() => {
    if (!statusData || !activeWorkflowId) return;

    const currentStep = statusData.current_step;
    let newStatus: WorkflowStatus = "running";

    if (currentStep === "End") {
      newStatus = "completed";
      // 完成後停止輪詢
      setTimeout(() => setActiveWorkflowId(null), 2000);
    } else if (currentStep === "Failed") {
      newStatus = "failed";
      setTimeout(() => setActiveWorkflowId(null), 2000);
    } else if (currentStep === "Paused") {
      newStatus = "paused";
    }

    setWorkflowStatuses((prev) => ({
      ...prev,
      [activeWorkflowId]: {
        status: newStatus,
        currentStep: newStatus === "running" ? currentStep : undefined,
      },
    }));
  }, [statusData, activeWorkflowId]);

  const handleTrigger = (workflowId: string) => {
    const workflow = workflows.find((w: WorkflowPayload) => w.workflow_id === workflowId);
    if (!workflow) return;

    triggerWorkflow.mutate(workflow, {
      onSuccess: () => {
        setActiveWorkflowId(workflowId);
        setWorkflowStatuses((prev) => ({
          ...prev,
          [workflowId]: { status: "running" },
        }));
      },
      onError: (error) => {
        alert(`執行失敗: ${error.message}`);
      },
    });
  };

  const handlePause = (workflowId: string) => {
    const currentStatus = workflowStatuses[workflowId]?.status;
    
    if (currentStatus === "paused") {
      // Resume
      resumeWorkflow.mutate(workflowId, {
        onSuccess: () => {
          setActiveWorkflowId(workflowId);
          setWorkflowStatuses((prev) => ({
            ...prev,
            [workflowId]: { status: "running" },
          }));
        },
        onError: (error) => {
          alert(`恢復失敗: ${error.message}`);
        },
      });
    } else {
      // Pause
      pauseWorkflow.mutate(workflowId, {
        onSuccess: () => {
          setWorkflowStatuses((prev) => ({
            ...prev,
            [workflowId]: { status: "paused" },
          }));
        },
        onError: (error) => {
          alert(`暫停失敗: ${error.message}`);
        },
      });
    }
  };

  const handleDelete = (workflowId: string) => {
    deleteWorkflow.mutate(workflowId, {
      onSuccess: () => {
        // 從狀態中移除
        setWorkflowStatuses((prev) => {
          const newStatuses = { ...prev };
          delete newStatuses[workflowId];
          return newStatuses;
        });
      },
      onError: (error) => {
        alert(`刪除失敗: ${error.message}`);
      },
    });
  };

  const workflowItems: Array<{ id: string; name: string; status: WorkflowStatus; currentStep?: string }> = useMemo(() => {
    return workflows.map((w: WorkflowPayload) => ({
      id: w.workflow_id || "",
      name: w.workflow_name || "Untitled",
      status: workflowStatuses[w.workflow_id || ""]?.status || "idle" as WorkflowStatus,
      currentStep: workflowStatuses[w.workflow_id || ""]?.currentStep,
    }));
  }, [workflows, workflowStatuses]);

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-12">
        <div className="text-muted-foreground">Loading workflows...</div>
      </div>
    );
  }

  if (workflowItems.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center py-16 text-center">
        <div className="rounded-full bg-muted p-4 mb-4">
          <Workflow className="h-8 w-8 text-muted-foreground" />
        </div>
        <h3 className="font-semibold text-lg mb-2">還沒有任何 Workflow</h3>
        <p className="text-muted-foreground text-sm">
          點擊上方的「Create Workflow」按鈕來建立你的第一個 Workflow
        </p>
      </div>
    );
  }

  return (
    <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
      {workflowItems.map((workflow) => (
        <WorkflowCard
          key={workflow.id}
          workflowId={workflow.id}
          workflowName={workflow.name}
          status={workflow.status}
          currentStep={workflow.currentStep}
          onTrigger={handleTrigger}
          onPause={handlePause}
          onDelete={handleDelete}
        />
      ))}
    </div>
  );
}
