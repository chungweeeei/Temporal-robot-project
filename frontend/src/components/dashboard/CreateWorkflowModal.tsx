import { useState } from "react";
import { useNavigate } from "react-router-dom";
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
import { useCreateWorkflow } from "@/hooks/useCreateWorkflow";
import { Plus } from "lucide-react";

export function CreateWorkflowModal() {
  const [open, setOpen] = useState(false);
  const [workflowName, setWorkflowName] = useState("");
  const navigate = useNavigate();
  const createWorkflow = useCreateWorkflow();

  const handleCreate = () => {
    if (!workflowName.trim()) {
      alert("請輸入 Workflow 名稱");
      return;
    }

    createWorkflow.mutate(workflowName.trim(), {
      onSuccess: (data) => {
        setOpen(false);
        setWorkflowName("");
        // 建立成功後導航至 Editor 編輯新的 workflow
        if (data?.workflow_id) {
          navigate(`/editor/${data.workflow_id}`);
        }
      },
      onError: (error) => {
        alert(`建立失敗: ${error.message}`);
      },
    });
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter") {
      handleCreate();
    }
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button>
          <Plus className="h-4 w-4 mr-2" />
          Create Workflow
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>建立新的 Workflow</DialogTitle>
          <DialogDescription>
            為你的新 Workflow 命名，建立後可以在編輯器中設計流程。
          </DialogDescription>
        </DialogHeader>
        <div className="py-4">
          <Input
            id="workflow-name"
            placeholder="輸入 Workflow 名稱..."
            value={workflowName}
            onChange={(e) => setWorkflowName(e.target.value)}
            onKeyDown={handleKeyDown}
            autoFocus
          />
        </div>
        <DialogFooter>
          <Button variant="outline" onClick={() => setOpen(false)}>
            取消
          </Button>
          <Button onClick={handleCreate} disabled={createWorkflow.isPending}>
            {createWorkflow.isPending ? "建立中..." : "建立"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
