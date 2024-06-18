## API Stability Levels
When a new API is introduced, you may want to allow developers to change its behavior without triggering a breaking change error.  
You can define an endpoint's stability level with the `x-stability-level` extension.  
There are four stability levels: `draft`->`alpha`->`beta`->`stable`.  
APIs with the levels `draft` or `alpha` can be changed freely without triggering a breaking change error.  
Stability level may be increased, but not decreased, like this: `draft`->`alpha`->`beta`->`stable`.  
APIs with no stability level will trigger breaking changes errors upon relevant changes.  
APIs with no stability level can be changed to any stability level.  

Example:
   ```
   /api/test:
    post:
     x-stability-level: "alpha"
   ```

Stability levels can also be used to control [grace periods for API deprecation](DEPRECATION.md#grace-period).