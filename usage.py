import asyncio
import subprocess

async def run_go_program():
    process = await asyncio.create_subprocess_exec(
        'go', 'run', 'main.go',
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE)

    # Handle standard output for regular data
    async for line in process.stdout:
        print(f'Received post: {line.decode().strip()}')

    # Handle standard error for logs
    async for line in process.stderr:
        print(f'Log output: {line.decode().strip()}')

    await process.wait()

async def main():
    await run_go_program()

if __name__ == '__main__':
    asyncio.run(main())
