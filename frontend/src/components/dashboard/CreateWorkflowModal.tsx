import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Plus } from "lucide-react";
import { CustomDialog } from "@/components/shared/CustomDialog";

export function CreateWorkflowModal() {

  // State for controlling the modal open state and workflow name input
  const [open, setOpen] = useState(false);
  const [workflowName, setWorkflowName] = useState("");
  
  // hook for react router navigation
  const navigate = useNavigate();


  return (
    <CustomDialog
      open={open}
      onOpenChange={setOpen}
      className="sm:max-w-[425px]"
      trigger={
        <Button>
          <Plus className="h-4 w-4 mr-2" />
          Create Workflow
        </Button>
      }
      title="Create New Workflow"
      description="Name your new workflow."
      footer={
        <>
          <Button variant="outline" onClick={() => setOpen(false)}>
            Cancel
          </Button>
          <Button onClick={() => {
            if (!workflowName.trim()) {
              alert("Please enter a workflow name.");
              return;
            }

            // Generate a new workflow ID 
            const newWorkflowId = `workflow-${Date.now()}`;

            // Navigate to the editor page with the new workflow ID
            navigate(`/editor/${newWorkflowId}`, {
              state: { 
                operation: "create",
                workflowName: workflowName.trim() 
              }
            });
          }}>
            Create
          </Button>
        </>
      }
    >
      <Input
        id="workflow-name"
        placeholder="Enter workflow name..."
        value={workflowName}
        onChange={(e) => setWorkflowName(e.target.value)}
        autoFocus
      />
    </CustomDialog>
  );
}
