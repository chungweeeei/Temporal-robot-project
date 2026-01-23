import { useState } from 'react';
import { CustomDialog } from '../shared/CustomDialog';
import { Sparkles, Loader2 } from "lucide-react";
import { Button } from "@/components/ui/button";

export function AIGenerateFlowModal() {
    
    const [open, setOpen] = useState(false);
    const [prompt, setPrompt] = useState("");

    const isLoading = false;

    return(
        <CustomDialog
            open={open}
            onOpenChange={setOpen}
            trigger={
              <Button variant="outline" size="sm">
                <Sparkles className="h-4 w-4 mr-2" />
                AI Generate
              </Button>
            }
            title="Generate Workflow with AI"
            description="Describe what you want the robot to do, and AI will generate the workflow for you."
            footer={
                <div className="flex gap-2">
                  <Button variant="outline" onClick={() => setOpen(false)}>
                    Cancel
                  </Button>
                  <Button onClick={() => console.log("click")} disabled={!prompt.trim() || isLoading}>
                    {isLoading ? (
                      <>
                        <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                        Generating...
                      </>
                    ) : (
                      <>
                        <Sparkles className="h-4 w-4 mr-2" />
                        Generate
                      </>
                    )}
                  </Button>
                </div>
            }
            >
            <div className="space-y-4">
              <textarea
                value={prompt}
                onChange={(e) => setPrompt(e.target.value)}
                placeholder="e.g., Move to position (1, 2), then say hello, and finally sit down."
                className="w-full min-h-[120px] p-3 border border-gray-300 rounded-md resize-none focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none"
              />
              <p className="text-xs text-muted-foreground">
                Tip: Be specific about the actions you want the robot to perform.
              </p>
            </div>
        </CustomDialog>
    )

}