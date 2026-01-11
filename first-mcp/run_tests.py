#!/usr/bin/env python3
"""
Test runner script for MCP Server
Provides convenient ways to run different types of tests
"""

import sys
import subprocess
import argparse
from pathlib import Path


def run_command(cmd, description):
    """Run a command and handle the result"""
    print(f"\n{'='*50}")
    print(f"Running: {description}")
    print(f"Command: {' '.join(cmd)}")
    print(f"{'='*50}")
    
    try:
        result = subprocess.run(cmd, check=True, capture_output=True, text=True)
        print(result.stdout)
        if result.stderr:
            print("STDERR:", result.stderr)
        return True
    except subprocess.CalledProcessError as e:
        print(f"❌ Failed: {description}")
        print(f"Exit code: {e.returncode}")
        print(f"STDOUT: {e.stdout}")
        print(f"STDERR: {e.stderr}")
        return False


def main():
    parser = argparse.ArgumentParser(description="Run MCP Server tests")
    parser.add_argument(
        "--type", 
        choices=["unit", "integration", "all"], 
        default="all",
        help="Type of tests to run"
    )
    parser.add_argument(
        "--verbose", "-v", 
        action="store_true",
        help="Verbose output"
    )
    parser.add_argument(
        "--coverage", 
        action="store_true",
        help="Run with coverage report"
    )
    
    args = parser.parse_args()
    
    # Check if pytest is available
    try:
        subprocess.run(["python", "-m", "pytest", "--version"], 
                      check=True, capture_output=True)
    except subprocess.CalledProcessError:
        print("❌ pytest not found. Please install it with: pip install pytest pytest-asyncio")
        return 1
    
    success = True
    
    if args.type in ["unit", "all"]:
        cmd = ["python", "-m", "pytest", "test_server.py"]
        if args.verbose:
            cmd.append("-v")
        if args.coverage:
            cmd.extend(["--cov=server", "--cov-report=term-missing"])
        
        if not run_command(cmd, "Unit Tests"):
            success = False
    
    if args.type in ["integration", "all"]:
        cmd = ["python", "-m", "pytest", "test_integration.py", "-m", "integration"]
        if args.verbose:
            cmd.append("-v")
        
        if not run_command(cmd, "Integration Tests"):
            success = False
    
    if args.type == "all":
        # Run all tests together for final summary
        cmd = ["python", "-m", "pytest"]
        if args.verbose:
            cmd.append("-v")
        if args.coverage:
            cmd.extend(["--cov=server", "--cov-report=term-missing", "--cov-report=html"])
        
        print(f"\n{'='*50}")
        print("Running All Tests Together")
        print(f"{'='*50}")
        
        if not run_command(cmd, "All Tests"):
            success = False
    
    # Summary
    print(f"\n{'='*50}")
    if success:
        print("✅ All tests passed!")
    else:
        print("❌ Some tests failed!")
    print(f"{'='*50}")
    
    return 0 if success else 1


if __name__ == "__main__":
    sys.exit(main())