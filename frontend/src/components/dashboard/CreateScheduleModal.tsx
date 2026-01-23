import { useState } from "react";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useCreateSchedule } from "@/hooks/useCreateSchedule";
import { useFetchWorkflow } from "@/hooks/useFetchWorkflows";
import { Plus } from "lucide-react";
import type { WorkflowInfo } from "@/types/workflows";

export function CreateScheduleModal() {
  const [open, setOpen] = useState(false);
  const [workflowId, setWorkflowId] = useState("");
  const [cronExpr, setCronExpr] = useState("*/5 * * * *");

  const createSchedule = useCreateSchedule();
  const { data: workflows = [] } = useFetchWorkflow();

  const handleCreate = () => {
    if (!workflowId || !cronExpr.trim()) {
      alert("Please fill in all required fields.");
      return;
    }

    createSchedule.mutate(
      {
        schedule_id: "schedule-" + Date.now(),
        workflow_id: workflowId,
        cron_expr: cronExpr.trim(),
      },
      {
        onSuccess: () => {
          setOpen(false);
          setWorkflowId("");
          setCronExpr("*/5 * * * *");
        },
        onError: (error) => {
          alert(`Failed to create schedule: ${error.message}`);
        },
      }
    );
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button variant="outline">
          <Plus className="h-4 w-4 mr-2" />
          Create Schedule
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Create New Schedule</DialogTitle>
          <DialogDescription>
            Create an automatic schedule for an existing Workflow.
          </DialogDescription>
        </DialogHeader>
        <div className="grid gap-4 py-4">
          <div className="space-y-2">
            <label className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
              Workflow
            </label>
            <select
              className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
              value={workflowId}
              onChange={(e) => setWorkflowId(e.target.value)}
            >
              <option value="" disabled>
                Select Workflow
              </option>
              {workflows.map((wf: WorkflowInfo) => (
                <option key={wf.workflow_id} value={wf.workflow_id}>
                  {wf.workflow_name}
                </option>
              ))}
            </select>
          </div>
          <div className="space-y-2">
            <label className="text-sm font-medium leading-none">
              Cron Expression
            </label>
            <Input
              placeholder="*/5 * * * *"
              value={cronExpr}
              onChange={(e) => setCronExpr(e.target.value)}
            />
          </div>
        </div>
        <DialogFooter>
          <Button variant="outline" onClick={() => setOpen(false)}>
            Cancel
          </Button>
          <Button onClick={handleCreate} disabled={createSchedule.isPending}>
            {createSchedule.isPending ? "Creating..." : "Create"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
