package main

const Source string = `
#include <uapi/linux/ptrace.h>
#include <linux/sched.h>
#include <linux/fs.h>

#define ARGSIZE 128

struct output_format
{
    u64 ts;
	char task[TASK_COMM_LEN];
    u32 pid;
    u32 type;
    char argv[ARGSIZE];
};

BPF_PERF_OUTPUT(events);

int syscall__execve(struct pt_regs *ctx,
    const char __user * filename,
    const char __user * const __user * __argv,
    const char __user * const __user * __envp)
{
	struct output_format data = {};
    struct task_struct *task = (struct task_struct *)bpf_get_current_task();

    data.ts = bpf_ktime_get_ns();
    bpf_get_current_comm(&data.task, sizeof(data.task));
    data.pid = bpf_get_current_pid_tgid();
    data.type = 1;
    
    for(u64 i = 0; i < 32; i++)
    {
        const char *argp = NULL;
        bpf_probe_read_user(&argp, sizeof(argp), (void *)&__argv[i]);
        if(argp)
        {
            bpf_probe_read_user(data.argv, sizeof(data.argv), argp);
            events.perf_submit(ctx, &data, sizeof(data));
        }
        else
            break;
    }

    data.type = 0;
    events.perf_submit(ctx, &data, sizeof(data));
	return 0;
}
`
