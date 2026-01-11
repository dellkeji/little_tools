#!/usr/bin/env python3
"""
MCP Server with version query tool and extensible architecture
"""

import asyncio
import logging
from typing import Any, Sequence

from mcp.server import Server
from mcp.server.models import InitializationOptions
from mcp.server.stdio import stdio_server
from mcp.types import (
    CallToolRequest,
    CallToolResult,
    ListToolsRequest,
    TextContent,
    Tool,
)

# Server version
SERVER_VERSION = "1.0.0"

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

# Create server instance
server = Server("my-mcp-server")


@server.list_tools()
async def handle_list_tools() -> list[Tool]:
    """
    List available tools.
    Each tool should be defined here for discoverability.
    """
    return [
        Tool(
            name="get_server_version",
            description="Get the current version of this MCP server",
            inputSchema={
                "type": "object",
                "properties": {},
                "required": [],
            },
        ),
        # Add more tools here as needed
        # Tool(
        #     name="example_tool",
        #     description="Example tool description",
        #     inputSchema={
        #         "type": "object",
        #         "properties": {
        #             "param": {
        #                 "type": "string",
        #                 "description": "Example parameter"
        #             }
        #         },
        #         "required": ["param"],
        #     },
        # ),
    ]


@server.call_tool()
async def handle_call_tool(name: str, arguments: dict[str, Any] | None) -> list[TextContent]:
    """
    Handle tool calls.
    This is where you implement the logic for each tool.
    """
    if name == "get_server_version":
        return await get_server_version(arguments or {})
    
    # Add more tool handlers here
    # elif name == "example_tool":
    #     return await example_tool_handler(arguments or {})
    
    else:
        raise ValueError(f"Unknown tool: {name}")


async def get_server_version(arguments: dict[str, Any]) -> list[TextContent]:
    """
    Get the current server version.
    """
    logger.info("Getting server version")
    
    return [
        TextContent(
            type="text",
            text=f"MCP Server Version: {SERVER_VERSION}"
        )
    ]


# Example of how to add a new tool handler
async def example_tool_handler(arguments: dict[str, Any]) -> list[TextContent]:
    """
    Example tool handler - uncomment and modify as needed.
    """
    param = arguments.get("param", "")
    
    # Your tool logic here
    result = f"Processed parameter: {param}"
    
    return [
        TextContent(
            type="text",
            text=result
        )
    ]


async def main():
    """
    Main entry point for the server.
    """
    # Run the server using stdio transport
    async with stdio_server() as (read_stream, write_stream):
        await server.run(
            read_stream,
            write_stream,
            InitializationOptions(
                server_name="my-mcp-server",
                server_version=SERVER_VERSION,
                capabilities=server.get_capabilities(
                    notification_options=None,
                    experimental_capabilities=None,
                ),
            ),
        )


if __name__ == "__main__":
    asyncio.run(main())