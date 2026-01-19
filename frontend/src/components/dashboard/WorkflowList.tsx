import { useMemo, useEffect } from "react";
import { WorkflowCard } from "./WorkflowCard";
import { useFetchWorkflow } from "@/hooks/useFetchWorkflows";
import { useFetchWorkflowStatus } from "@/hooks/useFetchWorkflowStatus";
import { useTriggerWorkflow, usePauseWorkflow, useResumeWorkflow } from "@/hooks/useControlWorkflow";
import { useDeleteWorkflow } from "@/hooks/useDeleteWorkflow";
import { useWorkflowStore } from "@/store/useWorkflowStore";
import type { WorkflowInfo, WorkflowStatusDef } from "@/types/workflows";
import { Workflow } from "lucide-react";

export function WorkflowList() {
  const { data: workflows = [], isLoading } = useFetchWorkflow();
  const { activeWorkflowId, setActiveWorkflowId } = useWorkflowStore();

  // fetch status for the active workflow
  const { data: activeWorkflowStatus } = useFetchWorkflowStatus(
    activeWorkflowId || "", 
    !!activeWorkflowId
  );

  // Register mutation hooks
  const triggerWorkflow = useTriggerWorkflow();
  const pauseWorkflow = usePauseWorkflow();
  const resumeWorkflow = useResumeWorkflow();
  const deleteWorkflow = useDeleteWorkflow();

  // derive workflow statuses from activeWorkflowStatus
  const workflowStatuses = useMemo<Record<string, { status: WorkflowStatusDef; current_step: string }>>(() => {
    if (!activeWorkflowId || !activeWorkflowStatus) return {};
    return {
      [activeWorkflowId]: {
        status: activeWorkflowStatus.status,
        current_step: activeWorkflowStatus.current_step || "",
      },
    };
  }, [activeWorkflowId, activeWorkflowStatus]);

  // Handle side effects for completed/failed workflows
  useEffect(() => {
    if (!activeWorkflowStatus || !activeWorkflowId) return;

    if (activeWorkflowStatus.status === "Completed" || activeWorkflowStatus.status === "Failed") {
      const timer = setTimeout(() => setActiveWorkflowId(null), 3000);
      return () => clearTimeout(timer);
    }
  }, [activeWorkflowId, activeWorkflowStatus, setActiveWorkflowId]);

  const handleTrigger = (workflowId: string) => {
    const workflow = workflows.find((w: WorkflowInfo) => w.workflow_id === workflowId);
    if (!workflow) return;

    triggerWorkflow.mutate(workflow, {
      onSuccess: () => {
        setActiveWorkflowId(workflowId);
      },
      onError: (error) => {
        alert(`Failed to trigger workflow: ${error.message}`);
      },
    });
  };

  const handlePause = (workflowId: string) => {
    const currentStatus = workflowStatuses[workflowId]?.status;

    if (currentStatus === "Paused") {
      // Resume
      resumeWorkflow.mutate(workflowId, {
        onSuccess: () => {
          setActiveWorkflowId(workflowId);
        },
        onError: (error) => {
          alert(`Failed to resume: ${error.message}`);
        },
      });
    } else {
      // Pause
      pauseWorkflow.mutate(workflowId, {
        onSuccess: () => {
          console.log("Workflow paused successfully.");
        },
        onError: (error) => {
          alert(`Failed to pause: ${error.message}`);
        },
      });
    }
  };

  const handleDelete = (workflowId: string) => {
    deleteWorkflow.mutate(workflowId, {
      onSuccess: () => {
        console.log(`Workflow ${workflowId} deleted successfully.`);
      },
      onError: (error) => {
        alert(`Workflow ${workflowId} deletion failed: ${error.message}`);
      },
    });
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-12">
        <div className="text-muted-foreground">Loading workflows...</div>
      </div>
    );
  }

  if (workflows.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center py-16 text-center">
        <div className="rounded-full bg-muted p-4 mb-4">
          <Workflow className="h-8 w-8 text-muted-foreground" />
        </div>
        <h3 className="font-semibold text-lg mb-2">No Workflows Yet</h3>
        <p className="text-muted-foreground text-sm">
          Click above「Create Workflow」button to get started.
        </p>
      </div>
    );
  }
 
  return (
    <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
      {workflows.map((workflow) => (
        <WorkflowCard
          key={workflow.workflow_id}
          workflowId={workflow.workflow_id}
          workflowName={workflow.workflow_name}
          status={workflowStatuses[workflow.workflow_id]?.status || "Idle"}
          currentStep={workflowStatuses[workflow.workflow_id]?.current_step || ""}
          onDelete={handleDelete}
          onTrigger={handleTrigger}
          onPause={handlePause}
        />
      ))}
    </div>
  );
}
