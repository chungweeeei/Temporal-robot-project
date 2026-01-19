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

export function CreateScheduleModal() {
  const [open, setOpen] = useState(false);
  const [scheduleId, setScheduleId] = useState("");
  const [workflowId, setWorkflowId] = useState("");
  const [cronExpr, setCronExpr] = useState("*/5 * * * *");
  const [timezone, setTimezone] = useState("Asia/Taipei");

  const createSchedule = useCreateSchedule();
  const { data: workflows = [] } = useFetchWorkflow();

  const handleCreate = () => {
    if (!scheduleId.trim() || !workflowId || !cronExpr.trim()) {
      alert("請填寫完整資訊 (Timezone 預設 Asia/Taipei)");
      return;
    }

    createSchedule.mutate(
      {
        schedule_id: scheduleId.trim(),
        workflow_id: "RobotWorkflow",
        cron_expr: cronExpr.trim(),
        timezone: timezone.trim() || undefined,
      },
      {
        onSuccess: () => {
          setOpen(false);
          setScheduleId("");
          setWorkflowId("");
          // Reset to defaults if desired
        },
        onError: (error) => {
          alert(`建立 Schedule 失敗: ${error.message}`);
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
          <DialogTitle>建立新的 Schedule</DialogTitle>
          <DialogDescription>
            為現有的 Workflow 建立自動排程。
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
                選擇 Workflow
              </option>
              {workflows.map((wf: any) => (
                <option key={wf.workflow_id} value={wf.workflow_id}>
                  {wf.workflow_name || wf.workflow_id}
                </option>
              ))}
            </select>
          </div>
          <div className="space-y-2">
            <label className="text-sm font-medium leading-none">
              Schedule ID
            </label>
            <Input
              placeholder="e.g. task-check-schedule-001"
              value={scheduleId}
              onChange={(e) => setScheduleId(e.target.value)}
            />
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
          <div className="space-y-2">
            <label className="text-sm font-medium leading-none">
              Timezone
            </label>
            <Input
              placeholder="Asia/Taipei"
              value={timezone}
              onChange={(e) => setTimezone(e.target.value)}
            />
          </div>
        </div>
        <DialogFooter>
          <Button variant="outline" onClick={() => setOpen(false)}>
            取消
          </Button>
          <Button onClick={handleCreate} disabled={createSchedule.isPending}>
            {createSchedule.isPending ? "建立中..." : "建立"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
