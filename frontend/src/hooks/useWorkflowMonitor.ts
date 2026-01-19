import { useState, useEffect, useMemo } from "react";
import { useFetchWorkflowStatus } from './useFetchWorkflowStatus';
import { type WorkflowStatus } from '../types/schema';


export function useWorkflowMonitor(workflowId: string){

    const [isMonitoring, setIsMonitoring] = useState(false);

    const { data: statusData } = useFetchWorkflowStatus(workflowId || '', isMonitoring);

    const workflowStatus: WorkflowStatus = useMemo(() => {
        if (!statusData) return isMonitoring ? 'running' : 'idle';
    
        const step = statusData.current_step;
        if (step === 'End') return 'completed';
        if (step === 'Failed') return 'failed';
        if (step === 'Paused') return 'paused';
    
        return isMonitoring ? 'running' : 'idle';
    }, [statusData, isMonitoring])

    const currentStep = statusData?.current_step || null;

    // Effect: Auto-stop monitoring
    useEffect(() => {
      if (workflowStatus === 'completed' || workflowStatus === 'failed') {
        const timer = setTimeout(() => setIsMonitoring(false), 2000);
        return () => clearTimeout(timer);
      }
    }, [workflowStatus]);

    return {
      isMonitoring,
      setIsMonitoring,
      workflowStatus,
      currentStep
    };
}