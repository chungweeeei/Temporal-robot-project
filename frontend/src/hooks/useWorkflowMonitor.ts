import { useEffect, useMemo } from "react";
import { useFetchWorkflowStatus } from './useFetchWorkflowStatus';
import { type WorkflowStatus } from '../types/schema';
import { useWorkflowStore } from '../store/useWorkflowStore';

export function useWorkflowMonitor(workflowId: string) {
    const { activeWorkflowId, setActiveWorkflowId } = useWorkflowStore();
    
    // Determine if we are monitoring THIS specific workflow based on global state
    const isMonitoring = activeWorkflowId === workflowId;

    const { data: statusData } = useFetchWorkflowStatus(workflowId || '', isMonitoring);

    const workflowStatus: WorkflowStatus = useMemo(() => {
        if (!statusData) return isMonitoring ? 'running' : 'idle';
    
        const step = statusData.current_step;
        if (step === 'End') return 'completed';
        if (step === 'Failed') return 'failed';
        if (step === 'Paused') return 'paused';
    
        return isMonitoring ? 'running' : 'idle';
    }, [statusData, isMonitoring]);

    const currentStep = statusData?.current_step || null;

    // Effect: Auto-stop monitoring (Global)
    useEffect(() => {
      // Only affect global state if WE are the active workflow
      if (!isMonitoring) return;

      if (workflowStatus === 'completed' || workflowStatus === 'failed') {
        const timer = setTimeout(() => {
            // Check again if we are still the active one before clearing
            if (useWorkflowStore.getState().activeWorkflowId === workflowId) {
                setActiveWorkflowId(null);
            }
        }, 2000);
        return () => clearTimeout(timer);
      }
    }, [workflowStatus, isMonitoring, workflowId, setActiveWorkflowId]);

    // Helper to manually start/stop monitoring this workflow
    const setMonitoring = (enable: boolean) => {
        if (enable) {
            setActiveWorkflowId(workflowId);
        } else if (isMonitoring) {
            setActiveWorkflowId(null);
        }
    };

    return {
      isMonitoring,
      setIsMonitoring: setMonitoring,
      workflowStatus,
      currentStep
    };
}